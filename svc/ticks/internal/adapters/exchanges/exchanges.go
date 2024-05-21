package exchanges

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/pkg/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/svc/ticks/internal/adapters/exchanges/binance"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
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

func (e Exchanges) ListenSymbol(ctx context.Context, sub event.TickSubscription) (chan tick.Tick, chan struct{}, error) {
	switch sub.Exchange {
	case binancePkg.Infos.Name:
		return e.binance.ListenSymbol(ctx, sub.Pair)
	default:
		return nil, nil, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, sub.Exchange)
	}
}
