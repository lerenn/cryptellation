package backtests

import (
	"context"
	"fmt"
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

// Test interface implementation
var _ Operator = (*Backtests)(nil)

type Backtests struct {
	repository vdb.Port
	pubsub     pubsub.Port
	csClient   candlesticks.Client
}

func New(repository vdb.Port, ps pubsub.Port, csClient candlesticks.Client) *Backtests {
	if repository == nil {
		panic("nil repository")
	}

	if ps == nil {
		panic("nil pubsub")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &Backtests{
		repository: repository,
		pubsub:     ps,
		csClient:   csClient,
	}
}

func (b Backtests) Advance(ctx context.Context, backtestId uint) error {
	return b.repository.LockedBacktest(backtestId, func() error {
		// Get backtest info
		bt, err := b.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		// Advance backtest
		finished := bt.Advance()

		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = b.readActualEvents(ctx, bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				log.Println("WARNING: no event detected for", bt.CurrentCsTick.Time)
				bt.SetCurrentTime(bt.EndTime)
			} else if !evts[0].Time.Equal(bt.CurrentCsTick.Time) {
				log.Println("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].Time)
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		evts = append(evts, event.NewStatusEvent(bt.CurrentCsTick.Time, status.Status{
			Finished: finished,
		}))
		b.broadcastEvents(backtestId, evts)

		if len(evts) > 1 {
			if err := b.repository.UpdateBacktest(ctx, bt); err != nil {
				return fmt.Errorf("cannot update backtest: %w", err)
			}
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

		t, cs, exists := list.First()
		if !exists {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.ExchangeName, sub.PairSymbol, bt.CurrentCsTick.PriceType, t, cs)
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

func (b Backtests) CreateOrder(ctx context.Context, backtestId uint, order order.Order) error {
	return b.repository.LockedBacktest(backtestId, func() error {
		bt, err := b.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		list, err := b.csClient.ReadCandlesticks(ctx, candlesticks.ReadCandlestickPayload{
			ExchangeName: order.ExchangeName,
			PairSymbol:   order.PairSymbol,
			Period:       bt.PeriodBetweenEvents,
			Start:        bt.CurrentCsTick.Time,
			End:          bt.CurrentCsTick.Time,
			Limit:        0,
		})
		if err != nil {
			return fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		_, cs, notEmpty := list.First()
		if !notEmpty {
			return backtest.ErrNoDataForOrderValidation
		}

		if err := bt.AddOrder(order, cs); err != nil {
			return err
		}

		if err := b.repository.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}

func (b Backtests) Create(ctx context.Context, req backtest.NewPayload) (id uint, err error) {
	bt, err := backtest.New(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("creating a new backtest from request: %w", err)
	}

	err = b.repository.CreateBacktest(ctx, &bt)
	if err != nil {
		return 0, fmt.Errorf("adding backtest to vdb: %w", err)
	}

	return bt.ID, nil
}

func (b Backtests) GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error) {
	bt, err := b.repository.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}

func (b Backtests) GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error) {
	bt, err := b.repository.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error {
	return b.repository.LockedBacktest(backtestId, func() error {
		bt, err := b.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		if _, err = bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		if err := b.repository.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
