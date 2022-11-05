package service

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type MockExchangeService struct {
}

func (mes *MockExchangeService) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	})

	for i := payload.Start.Unix(); i < 60*1000; i += 60 {
		if time.Unix(i, 0).After(payload.End) {
			break
		}

		if err := cl.Set(time.Unix(i, 0), candlestick.Candlestick{
			Open:  float64(i),
			Close: 1234,
		}); err != nil {
			return nil, err
		}
	}

	return cl, nil
}
