package exchanges

import (
	"context"
	"fmt"

	binancePkg "cryptellation/internal/adapters/exchanges/binance"
	"cryptellation/pkg/config"

	"cryptellation/svc/candlesticks/internal/adapters/exchanges/binance"
	"cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"cryptellation/svc/candlesticks/pkg/candlestick"
)

type Exchanges struct {
	binance *binance.Service
}

func New() (Exchanges, error) {
	b, err := binance.New(config.LoadBinanceTest())
	if err != nil {
		return Exchanges{}, err
	}

	return Exchanges{
		binance: b,
	}, nil
}

func (e Exchanges) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	switch payload.Exchange {
	case binancePkg.Infos.Name:
		return e.binance.GetCandlesticks(ctx, payload)
	default:
		return nil, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, payload.Exchange)
	}
}
