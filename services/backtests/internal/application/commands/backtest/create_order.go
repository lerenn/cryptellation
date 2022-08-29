package cmdBacktest

import (
	"context"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/domain/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	candlesticksProto "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
)

type CreateOrderHandler struct {
	repository vdb.Port
	csClient   candlesticksProto.CandlesticksServiceClient
}

func NewCreateOrderHandler(repository vdb.Port, csClient candlesticksProto.CandlesticksServiceClient) CreateOrderHandler {
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

		resp, err := h.csClient.ReadCandlesticks(ctx, &candlesticksProto.ReadCandlesticksRequest{
			ExchangeName: order.ExchangeName,
			PairSymbol:   order.PairSymbol,
			PeriodSymbol: bt.PeriodBetweenEvents.String(),
			Start:        bt.CurrentCsTick.Time.Format(time.RFC3339),
			End:          bt.CurrentCsTick.Time.Format(time.RFC3339),
			Limit:        0,
		})
		if err != nil {
			return fmt.Errorf("could not get candlesticks from service: %w", err)
		} else if len(resp.Candlesticks) == 0 {
			return backtest.ErrNoDataForOrderValidation
		}

		if err := bt.AddOrder(order, resp.Candlesticks[0]); err != nil {
			return err
		}

		if err := h.repository.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
