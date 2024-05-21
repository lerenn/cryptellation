package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func (app Candlesticks) GetCached(ctx context.Context, payload app.GetCachedPayload) (*candlestick.List, error) {
	telemetry.L(ctx).Infof("Requests candlesticks from %s to %s (limit: %d)", payload.Start, payload.End, payload.Limit)

	// Be sure that we do not try to get data in the future
	fmt.Println("payload.End", payload.End)
	if payload.End == nil || payload.End.After(time.Now()) {
		payload.End = utils.ToReference(time.Now())
	}

	start, end := payload.Period.RoundInterval(payload.Start, payload.End)
	cl := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)

	// Read candlesticks from database
	if err := app.db.ReadCandlesticks(ctx, cl, start, end, payload.Limit); err != nil {
		return nil, err
	}
	telemetry.L(ctx).Debugf("Read DB for %d candlesticks from %s to %s (limit: %d)", cl.Len(), start, end, payload.Limit)

	missingRanges := cl.GetMissingRange(start, end, payload.Limit)
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

	downloadStart, downloadEnd := getDownloadStartEndTimes(ranges, payload.Period)
	if err := app.download(ctx, cl, downloadStart, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	// Upsert candlesticks to DB
	if err := app.upsert(ctx, cl); err != nil {
		return nil, err
	}

	rl := cl.Extract(start, end, payload.Limit)

	return rl, nil
}

// getDownloadStartEndTimes gives start and end time for download
// Time order: start < end
func getDownloadStartEndTimes(ranges []timeserie.TimeRange, p period.Symbol) (time.Time, time.Time) {
	start, end := ranges[0].Start, ranges[len(ranges)-1].End
	count := end.Sub(start) / p.Duration()

	if count < MinimalRetrievedMissingCandlesticks {
		difference := MinimalRetrievedMissingCandlesticks - count
		start = start.Add(-p.Duration() * difference / 2)
		end = end.Add(p.Duration() * difference / 2)
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
			return false, csToInsert.Set(ts, cs)
		} else if !rcs.Equal(cs) {
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
