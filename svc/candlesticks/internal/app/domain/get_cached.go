package domain

import (
	"context"
	"errors"
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"

	"github.com/lerenn/cryptellation/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/candlesticks/pkg/period"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func (app Candlesticks) GetCached(ctx context.Context, payload app.GetCachedPayload) (*candlestick.List, error) {
	cl := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)
	telemetry.L(ctx).Infof(
		"Requests candlesticks from %s to %s (limit: %d)",
		payload.Start, payload.End, payload.Limit)

	// Check there is an end and that is not in the future
	if payload.End == nil || payload.End.After(time.Now()) {
		telemetry.L(ctx).Debug("End is not set or is in the future, setting it to now()")
		payload.End = utils.ToReference(time.Now())
	}

	// Check there is a start and that is before end
	if payload.Start == nil || payload.Start.After(*payload.End) {
		telemetry.L(ctx).Debug("Start is not set or is after end, setting it to end - period * MinimalRetrievedMissingCandlesticks")
		payload.Start = utils.ToReference(payload.End.Add(-payload.Period.Duration() * MinimalRetrievedMissingCandlesticks))
	}

	// Round down payload start and end
	payload.Start = utils.ToReference(payload.Period.RoundTime(*payload.Start))
	payload.End = utils.ToReference(payload.Period.RoundTime(*payload.End))

	// Read candlesticks from database
	if err := app.db.ReadCandlesticks(ctx, cl, *payload.Start, *payload.End, payload.Limit); err != nil {
		return nil, err
	}
	telemetry.L(ctx).Debugf(
		"Read DB for %d candlesticks from %s to %s (limit: %d)",
		cl.Len(), *payload.Start, *payload.End, payload.Limit)

	missingRanges := cl.GetMissingRange(*payload.Start, *payload.End, payload.Limit)
	uncompleteRanges := cl.GetUncompleteRange()
	ranges, err := timeserie.MergeTimeRanges(missingRanges, uncompleteRanges)
	if err != nil {
		return cl, err
	}

	if len(ranges) == 0 {
		telemetry.L(ctx).Infof("No candlestick missing, returning the list with %d candlesticks.", cl.Len())
		return cl, nil
	}
	telemetry.L(ctx).Debugf("Candlesticks are missing from DB: %+v", timeserie.TimeRangesToString(ranges))

	// Download missing candlesticks
	downloadStart, downloadEnd := getDownloadStartEndTimes(ranges, payload.Period)
	if err := app.download(ctx, cl, downloadStart, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	// Upsert candlesticks to DB
	if err := app.upsert(ctx, cl); err != nil {
		return nil, err
	}

	// Only return the requested candlesticks
	rl := cl.Extract(*payload.Start, *payload.End, payload.Limit)
	telemetry.L(ctx).Debugf("Returning %d candlesticks from %s to %s", rl.Len(), *payload.Start, *payload.End)

	return rl, nil
}

// getDownloadStartEndTimes gives start and end time for download
// Time order: start < end
func getDownloadStartEndTimes(ranges []timeserie.TimeRange, p period.Symbol) (time.Time, time.Time) {
	start, end := ranges[0].Start, ranges[len(ranges)-1].End
	count := end.Sub(start) / p.Duration()

	// If there is less than MinimalRetrievedMissingCandlesticks candlesticks to retrieve
	if count < MinimalRetrievedMissingCandlesticks {
		difference := MinimalRetrievedMissingCandlesticks - count
		start = start.Add(-p.Duration() * difference / 2)
		end = end.Add(p.Duration() * difference / 2)
	}

	// Check that end is not in the future
	if end.After(time.Now()) {
		end = time.Now()
	}

	return p.RoundInterval(&start, &end)
}

func (app Candlesticks) download(ctx context.Context, cl *candlestick.List, start, end time.Time, limit uint) error {
	payload := exchanges.GetCandlesticksPayload{
		Exchange: cl.Exchange,
		Pair:     cl.Pair,
		Period:   cl.Period,
		Start:    start,
		End:      end,
	}

	for {
		ncl, err := app.exchanges.GetCandlesticks(ctx, payload)
		if err != nil {
			return err
		}

		telemetry.L(ctx).Debugf(
			"Read exchange for %d candlesticks from %s to %s (limit: %d)",
			ncl.Len(), payload.Start, payload.End, payload.Limit,
		)

		if err := cl.Merge(ncl, nil); err != nil {
			return err
		}

		cl.ReplaceUncomplete(ncl)

		t, _, exists := ncl.Last()
		if !exists || !t.Before(end) || (limit != 0 && cl.Len() >= int(limit)) {
			break
		}

		payload.Start = t.Add(cl.Period.Duration())
	}

	// Fill missing candlesticks to let know that there is no more data on exchange
	return cl.FillMissing(start, end, candlestick.Candlestick{})
}

func (app Candlesticks) upsert(ctx context.Context, cl *candlestick.List) error {
	tStart, _, startExists := cl.First()
	tEnd, _, endExists := cl.Last()
	if !startExists || !endExists {
		return nil
	}

	rcl := candlestick.NewListFrom(cl)
	if err := app.db.ReadCandlesticks(ctx, rcl, tStart, tEnd, 0); err != nil {
		return err
	}

	csToInsert := candlestick.NewListFrom(cl)
	csToUpdate := candlestick.NewListFrom(cl)
	if err := cl.Loop(func(ts time.Time, cs candlestick.Candlestick) (bool, error) {
		rcs, exists := rcl.Get(ts)
		if !exists {
			telemetry.L(ctx).Debugf("Inserting candlestick %s with %+v", ts, cs)
			return false, csToInsert.Set(ts, cs)
		} else if !rcs.Equal(cs) {
			telemetry.L(ctx).Debugf("Updating candlestick %s with %+v", ts, cs)
			return false, csToUpdate.Set(ts, cs)
		}
		return false, nil
	}); err != nil {
		return err
	}

	var insertErr error
	if csToInsert.Len() > 0 {
		insertErr = app.db.CreateCandlesticks(ctx, csToInsert)
	}

	var updateErr error
	if csToUpdate.Len() > 0 {
		updateErr = app.db.UpdateCandlesticks(ctx, csToUpdate)
	}

	return errors.Join(insertErr, updateErr)
}
