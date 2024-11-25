package live

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges/live/binance"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the live exchanges activities, regrouping all exchanges activities.
type Activities struct {
	binance *binance.Activities
}

// New will create a new live activities.
func New() (Activities, error) {
	b, err := binance.New()
	if err != nil {
		return Activities{}, err
	}

	return Activities{
		binance: b,
	}, nil
}

// Register will register the activities.
func (a Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.GetExchangeInfo,
		activity.RegisterOptions{Name: exchanges.GetExchangeInfoActivityName})
	w.RegisterActivityWithOptions(
		a.ListExchangesNames,
		activity.RegisterOptions{Name: exchanges.ListExchangesNamesActivityName})
}

// GetExchangeInfo will get a specific exchange info.
func (a Activities) GetExchangeInfo(
	ctx context.Context,
	params exchanges.GetExchangeInfoParams,
) (exchanges.GetExchangeInfoResult, error) {
	switch params.Name {
	case activities.BinanceInfos.Name:
		return a.binance.GetExchangeInfo(ctx, params)
	default:
		return exchanges.GetExchangeInfoResult{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Name)
	}
}

// ListExchangesNames will list the exchanges names.
func (a Activities) ListExchangesNames(
	ctx context.Context,
	params exchanges.ListExchangesNamesParams,
) (exchanges.ListExchangesNamesResult, error) {
	var names []string

	// Get Binance name
	binanceRes, err := a.binance.ListExchangesNames(ctx, params)
	if err != nil {
		return exchanges.ListExchangesNamesResult{}, err
	}
	names = append(names, binanceRes.List...)

	// Return the result
	return exchanges.ListExchangesNamesResult{
		List: names,
	}, nil
}
