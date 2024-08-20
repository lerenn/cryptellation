package exchanges

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/internal/adapters/exchanges/binance"

	"github.com/lerenn/cryptellation/svc/exchanges/internal/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
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
	case binancePkg.Infos.Name:
		return e.binance.Infos(ctx)
	default:
		return exchange.Exchange{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, name)
	}
}
