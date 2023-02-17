package exchangesAdapter

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges/binance"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/exchanges"
)

type Exchanges struct {
	binance *binance.Service
}

func New(c Config) (exchanges.Port, error) {
	b, err := binance.New(c.Binance)
	if err != nil {
		return nil, err
	}

	return Exchanges{
		binance: b,
	}, nil
}

func (e Exchanges) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	switch payload.Exchange {
	case binance.Name:
		return e.binance.GetCandlesticks(ctx, payload)
	default:
		return nil, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, payload.Exchange)
	}
}
