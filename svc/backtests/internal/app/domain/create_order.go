package domain

import (
	"context"
	"fmt"

	"cryptellation/internal/adapters/telemetry"
	"cryptellation/pkg/models/order"

	"cryptellation/svc/backtests/pkg/backtest"

	candlesticks "cryptellation/svc/candlesticks/clients/go"

	"github.com/google/uuid"
)

func (b Backtests) CreateOrder(ctx context.Context, backtestId uuid.UUID, order order.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

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

		telemetry.L(ctx).Infof("Adding %+v order to %q backtest", order, backtestId.String())
		if err := bt.AddOrder(order, cs); err != nil {
			return err
		}

		return nil
	})
}
