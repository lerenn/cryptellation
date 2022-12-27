package backtests

import (
	"context"
	"fmt"
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

func (b Backtests) Advance(ctx context.Context, backtestId uint) error {
	return b.db.LockedBacktest(backtestId, func() error {
		// Get backtest info
		bt, err := b.db.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		// Advance backtest
		finished := bt.Advance()

		// Get actual events
		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = b.readActualEvents(ctx, bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				log.Println("WARNING: no event detected for", bt.CurrentCsTick.Time)
				bt.SetCurrentTime(bt.EndTime)
				finished = true
			} else if !evts[0].Time.Equal(bt.CurrentCsTick.Time) {
				log.Println("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].Time)
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		// Add backtest status event
		evts = append(evts, event.NewStatusEvent(bt.CurrentCsTick.Time, status.Status{
			Finished: finished,
		}))
		b.broadcastEvents(backtestId, evts)

		// Update backtest
		if err := b.db.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}

func (b Backtests) readActualEvents(ctx context.Context, bt backtest.Backtest) ([]event.Event, error) {
	evts := make([]event.Event, 0, len(bt.TickSubscribers))
	for _, sub := range bt.TickSubscribers {
		list, err := b.csClient.ReadCandlesticks(ctx, candlesticks.ReadCandlestickPayload{
			ExchangeName: sub.ExchangeName,
			PairSymbol:   sub.PairSymbol,
			Period:       bt.PeriodBetweenEvents,
			Start:        bt.CurrentCsTick.Time,
			End:          bt.EndTime,
			Limit:        1,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		tcs, exists := list.First()
		if !exists {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.ExchangeName, sub.PairSymbol, bt.CurrentCsTick.PriceType, tcs.Time, tcs.Candlestick)
		if err != nil {
			return nil, fmt.Errorf("turning candlestick into event: %w", err)
		}
		evts = append(evts, evt)
	}

	_, evts = event.OnlyKeepEarliestSameTimeEvents(evts, bt.EndTime)
	return evts, nil
}

func (b Backtests) broadcastEvents(backtestId uint, evts []event.Event) {
	var count uint
	for _, evt := range evts {
		if err := b.pubsub.Publish(backtestId, evt); err != nil {
			log.Println("WARNING: error when publishing event", evt)
			continue
		}

		count++
	}

	if count == 0 {
		log.Println("WARNING: no available events")
	}
}
