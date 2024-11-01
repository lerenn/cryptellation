package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/event"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"

	"github.com/google/uuid"
)

func (b Backtests) Advance(ctx context.Context, backtestId uuid.UUID) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) (err error) {
		// Advance backtest
		finished, err := bt.Advance()
		if err != nil {
			return fmt.Errorf("cannot advance backtest: %w", err)
		}
		telemetry.L(ctx).Infof("Advancing backtest %s: %s", backtestId.String(), bt.CurrentTime())

		// Get actual events
		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = b.readActualEvents(ctx, *bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				telemetry.L(ctx).Info(fmt.Sprint("WARNING: no event detected for", bt.CurrentCandlestick.Time))
				bt.SetCurrentTime(bt.Parameters.EndTime)
				finished = true
			} else if !evts[0].Time.Equal(bt.CurrentCandlestick.Time) {
				telemetry.L(ctx).Info(fmt.Sprint("WARNING: no event between", bt.CurrentCandlestick.Time, "and", evts[0].Time))
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		// Add backtest status event
		evts = append(evts, event.NewStatusEvent(bt.CurrentCandlestick.Time, event.Status{
			Finished: finished,
		}))
		b.broadcastEvents(ctx, backtestId, evts)

		return nil
	})
}

func (b Backtests) readActualEvents(ctx context.Context, bt backtest.Backtest) ([]event.Event, error) {
	evts := make([]event.Event, 0, len(bt.PricesSubscriptions))
	for _, sub := range bt.PricesSubscriptions {
		list, err := b.candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
			Exchange: sub.Exchange,
			Pair:     sub.Pair,
			Period:   bt.Parameters.PricePeriod,
			Start:    &bt.CurrentCandlestick.Time,
			End:      &bt.Parameters.EndTime,
			Limit:    1,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		t, cs, exists := list.TimeSerie.First()
		if !exists {
			continue
		}

		evt, err := event.PriceEventFromCandlestick(sub.Exchange, sub.Pair, bt.CurrentCandlestick.Price, t, cs)
		if err != nil {
			return nil, fmt.Errorf("turning candlestick into event: %w", err)
		}
		evts = append(evts, evt)
	}

	_, evts = event.OnlyKeepEarliestSameTimeEvents(evts, bt.Parameters.EndTime)
	telemetry.L(ctx).Infof("%d events for ticks on backtest %s", len(evts), bt.ID.String())
	return evts, nil
}

func (b Backtests) broadcastEvents(ctx context.Context, backtestId uuid.UUID, evts []event.Event) {
	telemetry.L(ctx).Infof("Broadcasting %d events on backtest %s", len(evts), backtestId.String())

	var count uint
	for _, evt := range evts {
		telemetry.L(ctx).Infof("Broadcasting event %+v for backtest %s", evt, backtestId.String())
		if err := b.events.Publish(ctx, backtestId, evt); err != nil {
			telemetry.L(ctx).Info(fmt.Sprint("WARNING: error when publishing event", evt))
			continue
		}

		count++
	}

	if count == 0 {
		telemetry.L(ctx).Info("WARNING: no available events")
	}
}
