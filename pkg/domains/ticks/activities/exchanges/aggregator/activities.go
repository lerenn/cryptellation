package aggregator

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the struct that will handle the exchanges aggregator.
type Activities struct {
	exchanges map[string]exchanges.Exchanges
}

// New will create a new exchanges aggregator.
func New(exchs ...exchanges.Exchanges) exchanges.Exchanges {
	m := make(map[string]exchanges.Exchanges)
	for _, e := range exchs {
		m[e.Name()] = e
	}

	return &Activities{
		exchanges: m,
	}
}

// Name will return the name of the aggregator.
func (a *Activities) Name() string {
	return "aggregator"
}

// Register will register the exchanges aggregator to Temporal.
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
	exch, ok := a.exchanges[params.Exchange]
	if !ok {
		return exchanges.ListenSymbolResults{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
	}

	return exch.ListenSymbolActivity(ctx, params)
}
