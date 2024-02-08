package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func (app Candlesticks) GetCached(ctx context.Context, payload app.GetCachedPayload) (*candlestick.List, error) {
	telemetry.L(ctx).Info(fmt.Sprintf("Requests candlesticks from %s to %s (limit: %d)", payload.Start, payload.End, payload.Limit))

	// Be sure that we do not try to get data in the future
	if payload.End.After(time.Now()) {
		payload.End = utils.ToReference(time.Now())
	}

	start, end := payload.Period.RoundInterval(payload.Start, payload.End)
	cl := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)

	// Read candlesticks from database
	if err := app.db.ReadCandlesticks(ctx, cl, start, end, payload.Limit); err != nil {
		return nil, err
	}
	telemetry.L(ctx).Info(fmt.Sprintf("Read DB for %d candlesticks from %s to %s (limit: %d)", cl.Len(), start, end, payload.Limit))

	if !cl.AreMissing(start, end, payload.Limit) {
		telemetry.L(ctx).Info(fmt.Sprintf("No candlestick missing, returning the list with %d candlesticks.", cl.Len()))
		return cl, nil
	}
	telemetry.L(ctx).Info("Candlesticks are missing from DB")

	downloadStart, downloadEnd := getDownloadStartEndTimes(cl, start, end)
	if err := app.download(ctx, cl, downloadStart, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	if err := app.upsert(ctx, cl); err != nil {
		return nil, err
	}

	return cl.Extract(start, end, payload.Limit), nil
}

// getDownloadStartEndTimes gives start and end time for download
// Time order: start < end
func getDownloadStartEndTimes(cl *candlestick.List, end, start time.Time) (time.Time, time.Time) {
	t, _, exists := cl.TimeSerie.Last()
	if exists && !cl.HasUncomplete() {
		end = t.Add(cl.Period.Duration())
	}

	qty := int(cl.Period.CountBetweenTimes(end, start)) + 1
	if qty < MinimalRetrievedMissingCandlesticks {
		d := cl.Period.Duration() * time.Duration(MinimalRetrievedMissingCandlesticks-qty)
		start = start.Add(d)
	}

	return end, start
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

		msg := fmt.Sprintf(
			"Read exchange for %d candlesticks from %s to %s (limit: %d)",
			ncl.Len(), payload.Start, payload.End, payload.Limit)
		telemetry.L(ctx).Info(msg)

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
		_, exists := rcl.Get(ts)
		if !exists {
			return false, csToInsert.Set(ts, cs)
		} else {
			return false, csToUpdate.Set(ts, cs)
		}
	}); err != nil {
		return err
	}

	if csToInsert.Len() > 0 {
		return app.db.CreateCandlesticks(ctx, csToInsert)
	}

	if csToUpdate.Len() > 0 {
		return app.db.UpdateCandlesticks(ctx, csToUpdate)
	}

	return nil
}
