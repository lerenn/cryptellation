package aggregator

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
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

// Register will register the activities to Temporal.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.GetExchangeActivity,
		activity.RegisterOptions{Name: exchanges.GetExchangeActivityName})
	w.RegisterActivityWithOptions(
		a.ListExchangesActivity,
		activity.RegisterOptions{Name: exchanges.ListExchangesActivityName})
}

// GetExchangeActivity will get a specific exchange info.
func (a *Activities) GetExchangeActivity(
	ctx context.Context,
	params exchanges.GetExchangeActivityParams,
) (exchanges.GetExchangeActivityResults, error) {
	exch, ok := a.exchanges[params.Name]
	if !ok {
		return exchanges.GetExchangeActivityResults{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Name)
	}

	return exch.GetExchangeActivity(ctx, params)
}

// ListExchangesActivity will list the exchanges names.
func (a *Activities) ListExchangesActivity(
	ctx context.Context,
	params exchanges.ListExchangesActivityParams,
) (exchanges.ListExchangesActivityResults, error) {
	var names []string

	// Get Binance name
	for _, exch := range a.exchanges {
		binanceRes, err := exch.ListExchangesActivity(ctx, params)
		if err != nil {
			return exchanges.ListExchangesActivityResults{}, err
		}
		names = append(names, binanceRes.List...)
	}

	// Return the result
	return exchanges.ListExchangesActivityResults{
		List: names,
	}, nil
}
