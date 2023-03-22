package app

import (
	"context"
	"fmt"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/pkg/types/order"
)

func (b Backtests) CreateOrder(ctx context.Context, backtestId uint, order order.Order) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *domain.Backtest) error {
		list, err := b.candlesticks.Read(ctx, client.ReadCandlesticksPayload{
			ExchangeName: order.ExchangeName,
			PairSymbol:   order.PairSymbol,
			Period:       bt.PeriodBetweenEvents,
			Start:        &bt.CurrentCsTick.Time,
			End:          &bt.CurrentCsTick.Time,
			Limit:        0,
		})
		if err != nil {
			return fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		tcs, notEmpty := list.First()
		if !notEmpty {
			return domain.ErrNoDataForOrderValidation
		}

		if err := bt.AddOrder(order, tcs.Candlestick); err != nil {
			return err
		}

		return nil
	})
}
