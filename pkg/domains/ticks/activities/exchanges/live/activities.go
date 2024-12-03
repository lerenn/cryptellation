package live

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges/live/binance"
	"go.temporal.io/sdk/activity"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Activities is the struct that will handle the activities.
type Activities struct {
	binance *binance.Activities
}

// New will create a new exchanges activities.
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

// ListenSymbolActivity will listen to the symbol activity.
func (a Activities) ListenSymbolActivity(
	ctx context.Context,
	params exchanges.ListenSymbolParams,
) (exchanges.ListenSymbolResults, error) {
	switch params.Exchange {
	case binancePkg.BinanceInfos.Name:
		return a.binance.ListenSymbolActivity(ctx, params)
	default:
		return exchanges.ListenSymbolResults{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
	}
}
