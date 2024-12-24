package aggregator

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the live exchanges activities, regrouping all exchanges activities.
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

// Register will register the activities.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.GetCandlesticksActivity,
		activity.RegisterOptions{Name: exchanges.GetCandlesticksActivityName},
	)
}

// GetCandlesticksActivity will get the candlesticks.
func (a *Activities) GetCandlesticksActivity(
	ctx context.Context,
	params exchanges.GetCandlesticksActivityParams,
) (exchanges.GetCandlesticksActivityResults, error) {
	exch, ok := a.exchanges[params.Exchange]
	if !ok {
		return exchanges.GetCandlesticksActivityResults{},
			fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
	}

	return exch.GetCandlesticksActivity(ctx, params)
}
