package candlesticks

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/domain/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"golang.org/x/xerrors"
)

func (c Candlesticks) GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error) {
	start, end := candlesticks.ProcessRequestedStartEndTimes(payload.Period, payload.Start, payload.End)

	id := candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	}
	cl := candlestick.NewEmptyList(id)

	if err := c.db.ReadCandlesticks(ctx, cl, start, end, payload.Limit); err != nil {
		return nil, err
	}

	if !candlesticks.AreMissing(cl, start, end, payload.Limit) {
		return cl, nil
	}

	downloadStart, downloadEnd := candlesticks.GetDownloadStartEndTimes(cl, start, end)
	if err := c.download(ctx, cl, downloadStart, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	if err := c.upsert(ctx, cl); err != nil {
		return nil, err
	}

	return cl.Extract(start, end, payload.Limit), nil
}

func (reh Candlesticks) download(ctx context.Context, cl *candlestick.List, start, end time.Time, limit uint) error {
	exch, exists := reh.services[cl.ExchangeName()]
	if !exists {
		return xerrors.New(fmt.Sprintf("inexistant exchange service for %q", cl.ExchangeName()))
	}

	payload := exchanges.GetCandlesticksPayload{
		PairSymbol: cl.PairSymbol(),
		Period:     cl.Period(),
		Start:      start,
		End:        end,
	}

	for {
		ncl, err := exch.GetCandlesticks(ctx, payload)
		if err != nil {
			return err
		}

		if err := cl.Merge(*ncl, nil); err != nil {
			return err
		}

		if err := cl.ReplaceUncomplete(*ncl); err != nil {
			return err
		}

		c, exists := ncl.Last()
		if !exists || !c.Time.Before(end) || (limit != 0 && cl.Len() >= int(limit)) {
			break
		}

		payload.Start = c.Time.Add(cl.Period().Duration())
	}

	return nil
}

func (reh Candlesticks) upsert(ctx context.Context, cl *candlestick.List) error {
	start, startExists := cl.First()
	end, endExists := cl.Last()
	if !startExists || !endExists {
		return nil
	}

	rcl := candlestick.NewEmptyList(cl.ID())
	if err := reh.db.ReadCandlesticks(ctx, rcl, start.Time, end.Time, 0); err != nil {
		return err
	}

	csToInsert := candlestick.NewEmptyList(cl.ID())
	csToUpdate := candlestick.NewEmptyList(cl.ID())
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
		return reh.db.CreateCandlesticks(ctx, csToInsert)
	}

	if csToUpdate.Len() > 0 {
		return reh.db.UpdateCandlesticks(ctx, csToUpdate)
	}

	return nil
}
