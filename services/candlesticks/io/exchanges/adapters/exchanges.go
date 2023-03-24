package exchanges

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/services/candlesticks/io/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/io/exchanges/adapters/binance"
)

type Exchanges struct {
	binance *binance.Service
}

func New(c config.Exchanges) (Exchanges, error) {
	b, err := binance.New(c.Binance)
	if err != nil {
		return Exchanges{}, err
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
