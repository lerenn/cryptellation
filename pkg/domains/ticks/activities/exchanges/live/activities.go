package live

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.temporal.io/sdk/activity"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges/live/binance"
)

type Activities struct {
	binance *binance.Activities
}

func New(temporal temporalclient.Client) (exchanges.Exchanges, error) {
	// Load temporal configuration
	t := config.LoadTemporal(nil)
	if err := t.Validate(); err != nil {
		return nil, err
	}

	// Create binance adapter
	b, err := binance.New(temporal)
	if err != nil {
		return nil, err
	}

	return &Activities{
		binance: b,
	}, nil
}

// Register will register the activities.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.ListenSymbolActivity,
		activity.RegisterOptions{Name: exchanges.ListenSymbolActivityName})
}

func (e Activities) ListenSymbolActivity(ctx context.Context, params exchanges.ListenSymbolParams) (exchanges.ListenSymbolResults, error) {
	switch params.Exchange {
	case binancePkg.BinanceInfos.Name:
		return e.binance.ListenSymbol(ctx, params)
	default:
		return exchanges.ListenSymbolResults{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
	}
}
