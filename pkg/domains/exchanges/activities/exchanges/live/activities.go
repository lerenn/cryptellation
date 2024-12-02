package live

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges/live/binance"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the live exchanges activities, regrouping all exchanges activities.
type Activities struct {
	binance *binance.Activities
}

// New will create a new live activities.
func New() (*Activities, error) {
	b, err := binance.New()
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
	switch params.Name {
	case activities.BinanceInfos.Name:
		return a.binance.GetExchangeActivity(ctx, params)
	default:
		return exchanges.GetExchangeActivityResults{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Name)
	}
}

// ListExchangesActivity will list the exchanges names.
func (a *Activities) ListExchangesActivity(
	ctx context.Context,
	params exchanges.ListExchangesActivityParams,
) (exchanges.ListExchangesActivityResults, error) {
	var names []string

	// Get Binance name
	binanceRes, err := a.binance.ListExchangesActivity(ctx, params)
	if err != nil {
		return exchanges.ListExchangesActivityResults{}, err
	}
	names = append(names, binanceRes.List...)

	// Return the result
	return exchanges.ListExchangesActivityResults{
		List: names,
	}, nil
}
