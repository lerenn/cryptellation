package backtests

import (
	"context"
	"fmt"
	"log"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/backtest"
	"github.com/lerenn/cryptellation/pkg/models/order"
)

func (b Backtests) CreateOrder(ctx context.Context, backtestId uint, order order.Order) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
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
			return backtest.ErrNoDataForOrderValidation
		}

		log.Printf("Adding %+v order to %d backtest", order, backtestId)
		if err := bt.AddOrder(order, tcs.Candlestick); err != nil {
			return err
		}

		return nil
	})
}
