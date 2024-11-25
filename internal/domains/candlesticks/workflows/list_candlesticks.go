package workflows

import (
	"errors"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.temporal.io/sdk/workflow"
)

const (
	// MinimalRetrievedMissingCandlesticks is the minimal quantity of candlesticks
	// that will be retrieved in case of miss. It will avoid too many request on
	// exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

// ListCandlesticks is the workflow that will list candlesticks.
// TODO: Refactor this function
//
//nolint:funlen
func (c *candlesticks) ListCandlesticks(
	ctx workflow.Context,
	params api.ListCandlesticksParams,
) (api.ListCandlesticksResults, error) {
	// Log the start of the workflow
	workflow.GetLogger(ctx).Info(
		"Requested candlesticks started",
		"exchange", params.Exchange,
		"pair", params.Pair,
		"period", params.Period,
		"start", params.Start,
		"end", params.End,
		"limit", params.Limit)

	// Check and fix params
	params, err := validateCandlesticksParams(ctx, params)
	if err != nil {
		return api.ListCandlesticksResults{}, err
	}

	// Set activities params
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Read candlesticks from database
	var dbRes db.ReadCandlesticksResult
	err = workflow.ExecuteActivity(ctx, c.db.ReadCandlesticks, db.ReadCandlesticksParams{
		Exchange: params.Exchange,
		Pair:     params.Pair,
		Period:   params.Period,
		Start:    *params.Start,
		End:      *params.End,
		Limit:    params.Limit,
	}).Get(ctx, &dbRes)
	if err != nil {
		return api.ListCandlesticksResults{}, err
	}
	workflow.GetLogger(ctx).Debug(
		"DB candlesticks read executed",
		"retrieved", dbRes.List.Data.Len(),
		"from", *params.Start,
		"to", *params.End,
		"limit", params.Limit)

	// Spot missing candlesticks
	missingRanges := dbRes.List.GetMissingRange(*params.Start, *params.End, params.Limit)
	uncompleteRanges := dbRes.List.GetUncompleteRange()
	ranges, err := timeserie.MergeTimeRanges(missingRanges, uncompleteRanges)
	if err != nil {
		return api.ListCandlesticksResults{}, err
	}

	// If no candlesticks are missing, return the list
	if len(ranges) == 0 {
		workflow.GetLogger(ctx).Debug(
			"No candlestick missing, returning the candlesticks list.",
			"size", dbRes.List.Data.Len())
		return api.ListCandlesticksResults{
			List: dbRes.List,
		}, nil
	}
	workflow.GetLogger(ctx).Debug(
		"Candlesticks are missing from DB",
		"missing time ranges", timeserie.TimeRangesToString(ranges))

	// Download missing candlesticks
	downloadStart, downloadEnd := getDownloadStartEndTimes(ctx, ranges, params.Period)
	if err := c.download(ctx, dbRes.List, downloadStart, downloadEnd, params.Limit); err != nil {
		return api.ListCandlesticksResults{}, err
	}

	// Upsert candlesticks to DB
	if err := c.upsert(ctx, dbRes.List); err != nil {
		return api.ListCandlesticksResults{}, err
	}

	// Only return the requested candlesticks
	rl := dbRes.List.Extract(*params.Start, *params.End, params.Limit)
	workflow.GetLogger(ctx).Debug("Returning candlesticks to caller",
		"quantity", rl.Data.Len(),
		"from", *params.Start,
		"to", *params.End)

	return api.ListCandlesticksResults{List: rl}, nil
}

// getDownloadStartEndTimes gives start and end time for download.
// Time order: start < end.
func getDownloadStartEndTimes(
	ctx workflow.Context,
	ranges []timeserie.TimeRange,
	p period.Symbol,
) (time.Time, time.Time) {
	start, end := ranges[0].Start, ranges[len(ranges)-1].End
	count := end.Sub(start) / p.Duration()

	// If there is less than MinimalRetrievedMissingCandlesticks candlesticks to retrieve
	if count < MinimalRetrievedMissingCandlesticks {
		difference := MinimalRetrievedMissingCandlesticks - count
		start = start.Add(-p.Duration() * difference / 2)
		end = end.Add(p.Duration() * difference / 2)
	}

	// Check that end is not in the future
	if end.After(workflow.Now(ctx)) {
		end = workflow.Now(ctx)
	}

	return p.RoundInterval(&start, &end)
}

func (c candlesticks) download(ctx workflow.Context, cl *candlestick.List, start, end time.Time, limit uint) error {
	// Set params for download
	params := exchanges.GetCandlesticksParams{
		Exchange: cl.Metadata.Exchange,
		Pair:     cl.Metadata.Pair,
		Period:   cl.Metadata.Period,
		Start:    start,
		End:      end,
	}

	for {
		// Download candlesticks
		var exchangeRes exchanges.GetCandlesticksResult
		err := workflow.ExecuteActivity(ctx, c.exchanges.GetCandlesticks, params).Get(ctx, &exchangeRes)
		if err != nil {
			return err
		}
		workflow.GetLogger(ctx).Debug("Read exchange for candlesticks",
			"retrieved", exchangeRes.List.Data.Len(),
			"from", params.Start,
			"to", params.End,
			"limit", params.Limit)

		// Merge the downloaded candlesticks with the database candlesticks
		if err := cl.Merge(exchangeRes.List, nil); err != nil {
			return err
		}
		workflow.GetLogger(ctx).Debug("Merged candlesticks",
			"retrieved", exchangeRes.List.Data.Len(),
			"new size", cl.Data.Len())

		// Replace uncomplete candlesticks in the database candlesticks
		cl.ReplaceUncomplete(exchangeRes.List)

		// Check if there is no more data to retrieve
		t, _, exists := exchangeRes.List.Data.Last()
		if !exists || !t.Before(end) || (limit != 0 && cl.Data.Len() >= int(limit)) {
			break
		}

		// Set new start time for next download
		params.Start = t.Add(cl.Metadata.Period.Duration())
	}

	// Fill missing candlesticks to let know that there is no more data on exchange
	return cl.FillMissing(start, end, candlestick.Candlestick{})
}

func validateCandlesticksParams(
	ctx workflow.Context,
	params api.ListCandlesticksParams,
) (api.ListCandlesticksParams, error) {
	// Check if there is an exchange
	if params.Exchange == "" {
		return params, ErrNoExchange
	}

	// Check if there is a pair
	if params.Pair == "" {
		return params, ErrNoPair
	}

	// Check if there is a period
	if params.Period == "" {
		return params, ErrNoPeriod
	}

	// Check there is an end and that is not in the future
	if params.End == nil || params.End.After(workflow.Now(ctx)) {
		workflow.GetLogger(ctx).Debug("End is not set or is in the future, setting it to now()")
		params.End = utils.ToReference(workflow.Now(ctx))
	}

	// Check there is a start and that is before end
	if params.Start == nil || params.Start.After(*params.End) {
		workflow.GetLogger(ctx).Debug(
			"Start is not set or is after end, setting it to end - period * MinimalRetrievedMissingCandlesticks")
		params.Start = utils.ToReference(
			params.End.Add(-params.Period.Duration() * MinimalRetrievedMissingCandlesticks))
	}

	// Round down payload start and end
	params.Start = utils.ToReference(params.Period.RoundTime(*params.Start))
	params.End = utils.ToReference(params.Period.RoundTime(*params.End))

	return params, nil
}

// TODO: Refactor this function
//
//nolint:funlen,cyclop
func (c candlesticks) upsert(ctx workflow.Context, cl *candlestick.List) error {
	// Get the first and last candlestick
	// If there is no candlestick, return
	tStart, _, startExists := cl.Data.First()
	tEnd, _, endExists := cl.Data.Last()
	if !startExists || !endExists {
		return nil
	}

	// Read candlesticks from database
	var dbRes db.ReadCandlesticksResult
	err := workflow.ExecuteActivity(ctx, c.db.ReadCandlesticks, db.ReadCandlesticksParams{
		Exchange: cl.Metadata.Exchange,
		Pair:     cl.Metadata.Pair,
		Period:   cl.Metadata.Period,
		Start:    tStart,
		End:      tEnd,
	}).Get(ctx, &dbRes)
	if err != nil {
		return err
	}

	// Loop over the candlesticks and separate between to insert or update
	csToInsert := candlestick.NewListFrom(cl)
	csToUpdate := candlestick.NewListFrom(cl)
	if err := cl.Loop(func(cs candlestick.Candlestick) (bool, error) {
		rcs, exists := dbRes.List.Data.Get(cs.Time)
		if !exists {
			workflow.GetLogger(ctx).Debug("Inserting candlestick",
				"time", cs.Time,
				"data", cs)
			return false, csToInsert.Set(cs)
		} else if !rcs.Equal(cs) {
			workflow.GetLogger(ctx).Debug("Updating candlestick",
				"time", cs.Time,
				"data", cs)
			return false, csToUpdate.Set(cs)
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Insert candlesticks
	var insertErr error
	if csToInsert.Data.Len() > 0 {
		err := workflow.ExecuteActivity(ctx, c.db.CreateCandlesticks, db.CreateCandlesticksParams{
			List: csToInsert,
		}).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	// Update candlesticks
	var updateErr error
	if csToUpdate.Data.Len() > 0 {
		err := workflow.ExecuteActivity(ctx, c.db.UpdateCandlesticks, db.UpdateCandlesticksParams{
			List: csToUpdate,
		}).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return errors.Join(insertErr, updateErr)
}
