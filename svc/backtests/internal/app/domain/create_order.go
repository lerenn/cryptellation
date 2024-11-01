package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/order"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"

	"github.com/google/uuid"
)

func (b Backtests) CreateOrder(ctx context.Context, backtestId uuid.UUID, order order.Order) error {
	telemetry.L(ctx).Infof("Creating order on backtest %s: %+v", backtestId.String(), order)

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
		list, err := b.candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
			Exchange: order.Exchange,
			Pair:     order.Pair,
			Period:   bt.Parameters.PricePeriod,
			Start:    &bt.CurrentCandlestick.Time,
			End:      &bt.CurrentCandlestick.Time,
			Limit:    0,
		})
		if err != nil {
			return fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		_, cs, notEmpty := list.First()
		if !notEmpty {
			return fmt.Errorf("%w: %d candlesticks retrieved", backtest.ErrNoDataForOrderValidation, list.Len())
		}

		telemetry.L(ctx).Infof("Adding %+v order to %q backtest", order, backtestId.String())
		if err := bt.AddOrder(order, cs); err != nil {
			return err
		}

		return nil
	})
}
