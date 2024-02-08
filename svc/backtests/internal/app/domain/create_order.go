package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
)

func (b Backtests) CreateOrder(ctx context.Context, backtestId uint, order order.Order) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
		list, err := b.candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
			Exchange: order.Exchange,
			Pair:     order.Pair,
			Period:   bt.PeriodBetweenEvents,
			Start:    &bt.CurrentCsTick.Time,
			End:      &bt.CurrentCsTick.Time,
			Limit:    0,
		})
		if err != nil {
			return fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		_, cs, notEmpty := list.First()
		if !notEmpty {
			return backtest.ErrNoDataForOrderValidation
		}

		telemetry.L(ctx).Info(fmt.Sprintf("Adding %+v order to %d backtest", order, backtestId))
		if err := bt.AddOrder(order, cs); err != nil {
			return err
		}

		return nil
	})
}
