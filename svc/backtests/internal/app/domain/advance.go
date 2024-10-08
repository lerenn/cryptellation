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
		finished := bt.Advance()
		telemetry.L(ctx).Infof("Advancing backtest %d: %s", backtestId, bt.CurrentTime())

		// Get actual events
		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = b.readActualEvents(ctx, *bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				telemetry.L(ctx).Info(fmt.Sprint("WARNING: no event detected for", bt.CurrentCsTick.Time))
				bt.SetCurrentTime(bt.EndTime)
				finished = true
			} else if !evts[0].Time.Equal(bt.CurrentCsTick.Time) {
				telemetry.L(ctx).Info(fmt.Sprint("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].Time))
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		// Add backtest status event
		evts = append(evts, event.NewStatusEvent(bt.CurrentCsTick.Time, event.Status{
			Finished: finished,
		}))
		b.broadcastEvents(ctx, backtestId, evts)

		return nil
	})
}

func (b Backtests) readActualEvents(ctx context.Context, bt backtest.Backtest) ([]event.Event, error) {
	evts := make([]event.Event, 0, len(bt.TickSubscriptions))
	for _, sub := range bt.TickSubscriptions {
		list, err := b.candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
			Exchange: sub.Exchange,
			Pair:     sub.Pair,
			Period:   bt.PeriodBetweenEvents,
			Start:    &bt.CurrentCsTick.Time,
			End:      &bt.EndTime,
			Limit:    1,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		t, cs, exists := list.TimeSerie.First()
		if !exists {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.Exchange, sub.Pair, bt.CurrentCsTick.PriceType, t, cs)
		if err != nil {
			return nil, fmt.Errorf("turning candlestick into event: %w", err)
		}
		evts = append(evts, evt)
	}

	_, evts = event.OnlyKeepEarliestSameTimeEvents(evts, bt.EndTime)
	telemetry.L(ctx).Infof("%d events for ticks on backtest %d", len(evts), bt.ID)
	return evts, nil
}

func (b Backtests) broadcastEvents(ctx context.Context, backtestId uuid.UUID, evts []event.Event) {
	telemetry.L(ctx).Infof("Broadcasting %d events on backtest %d", len(evts), backtestId)

	var count uint
	for _, evt := range evts {
		telemetry.L(ctx).Infof("Broadcasting event %+v for backtest %d", evt, backtestId)
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
