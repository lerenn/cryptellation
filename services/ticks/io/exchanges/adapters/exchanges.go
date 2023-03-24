package exchanges

import (
	"fmt"

	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/digital-feather/cryptellation/services/ticks/io/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/io/exchanges/adapters/binance"
)

type Exchanges struct {
	binance *binance.Service
}

func New(c config.Exchanges) (Exchanges, error) {
	b, err := binance.New()
	if err != nil {
		return Exchanges{}, err
	}

	return Exchanges{
		binance: b,
	}, nil
}

func (e Exchanges) ListenSymbol(exchange, symbol string) (chan tick.Tick, chan struct{}, error) {
	switch exchange {
	case binance.Name:
		return e.binance.ListenSymbol(symbol)
	default:
		return nil, nil, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, exchange)
	}
}
