package candlesticks

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/domain/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"golang.org/x/xerrors"
)

// Test interface implementation
var _ Operator = (*Candlesticks)(nil)

type Candlesticks struct {
	repository db.Adapter
	services   map[string]exchanges.Adapter
}

func New(repository db.Adapter, services map[string]exchanges.Adapter) Candlesticks {
	if repository == nil {
		panic("nil repository")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return Candlesticks{
		repository: repository,
		services:   services,
	}
}

func (c Candlesticks) GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error) {
	start, end := candlesticks.ProcessRequestedStartEndTimes(payload.Period, payload.Start, payload.End)

	id := candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	}
	cl := candlestick.NewList(id)

	if err := c.repository.ReadCandlesticks(ctx, cl, start, end, payload.Limit); err != nil {
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

		t, _, exists := ncl.Last()
		if !exists || !t.Before(end) || (limit != 0 && cl.Len() >= int(limit)) {
			break
		}

		payload.Start = t.Add(cl.Period().Duration())
	}

	return nil
}

func (reh Candlesticks) upsert(ctx context.Context, cl *candlestick.List) error {
	start, _, startExists := cl.First()
	end, _, endExists := cl.Last()
	if !startExists || !endExists {
		return nil
	}

	rcl := candlestick.NewList(cl.ID())
	if err := reh.repository.ReadCandlesticks(ctx, rcl, start, end, 0); err != nil {
		return err
	}

	csToInsert := candlestick.NewList(cl.ID())
	csToUpdate := candlestick.NewList(cl.ID())
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
		return reh.repository.CreateCandlesticks(ctx, csToInsert)
	}

	if csToUpdate.Len() > 0 {
		return reh.repository.UpdateCandlesticks(ctx, csToUpdate)
	}

	return nil
}
