package exchanges

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/pkg/adapters/exchanges/port"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Exchanges struct {
	binance *binance.Service
}

func New() (Exchanges, error) {
	b, err := binance.New()
	if err != nil {
		return Exchanges{}, err
	}

	return Exchanges{
		binance: b,
	}, nil
}

func (e Exchanges) Infos(ctx context.Context, name string) (exchange.Exchange, error) {
	switch name {
	case binance.Infos.Name:
		return e.binance.Infos(ctx)
	default:
		return exchange.Exchange{}, fmt.Errorf("%w: %q", ErrInexistantExchange, name)
	}
}

func (e Exchanges) ListenSymbol(exchange, symbol string) (chan tick.Tick, chan struct{}, error) {
	switch exchange {
	case binance.Infos.Name:
		return e.binance.ListenSymbol(symbol)
	default:
		return nil, nil, fmt.Errorf("%w: %q", ErrInexistantExchange, exchange)
	}
}

func (e Exchanges) GetCandlesticks(ctx context.Context, payload port.GetCandlesticksPayload) (*candlestick.List, error) {
	switch payload.Exchange {
	case binance.Infos.Name:
		return e.binance.GetCandlesticks(ctx, payload)
	default:
		return nil, fmt.Errorf("%w: %q", ErrInexistantExchange, payload.Exchange)
	}
}
