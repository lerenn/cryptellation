package exchanges

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/internal/adapters/exchanges/binance"

	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/svc/candlesticks/internal/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
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
