package exchanges

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/exchange"
	"github.com/lerenn/cryptellation/services/exchanges/io/exchanges"
	"github.com/lerenn/cryptellation/services/exchanges/io/exchanges/adapters/binance"
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

func (e Exchanges) Infos(ctx context.Context, name string) (exchange.Exchange, error) {
	switch name {
	case binance.Infos.Name:
		return e.binance.Infos(ctx)
	default:
		return exchange.Exchange{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, name)
	}
}
