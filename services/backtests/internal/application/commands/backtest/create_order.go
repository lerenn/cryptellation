package cmdBacktest

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	candlesticksClient "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

type CreateOrderHandler struct {
	repository vdb.Port
	csClient   candlesticksClient.Client
}

func NewCreateOrderHandler(repository vdb.Port, csClient candlesticksClient.Client) CreateOrderHandler {
	if repository == nil {
		panic("nil repository")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return CreateOrderHandler{
		repository: repository,
		csClient:   csClient,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, backtestId uint, order order.Order) error {
	return h.repository.LockedBacktest(backtestId, func() error {
		bt, err := h.repository.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		list, err := h.csClient.ReadCandlesticks(ctx, candlesticksClient.ReadCandlestickPayload{
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

		if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
