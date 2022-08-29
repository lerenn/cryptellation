package service

import (
	"context"

	candlesticksClient "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type MockedCandlesticksClient struct {
}

func (m MockedCandlesticksClient) ReadCandlesticks(
	ctx context.Context,
	payload candlesticksClient.ReadCandlestickPayload,
) (*candlestick.List, error) {
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	})

	err := cl.Set(payload.Start, candlestick.Candlestick{
		Open:   1,
		High:   2,
		Low:    0.5,
		Close:  1.5,
		Volume: 500,
	})

	return cl, err
}
