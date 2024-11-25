package live

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"

	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges/live/binance"
)

type Activities struct {
	binance *binance.Activities
}

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
	w.RegisterActivityWithOptions(a.GetExchangeInfo, activity.RegisterOptions{Name: exchanges.GetExchangeInfoActivityName})
}

func (a Activities) GetExchangeInfo(ctx context.Context, params exchanges.GetExchangeInfoParams) (exchanges.GetExchangeInfoResult, error) {
	switch params.Name {
	case activities.BinanceInfos.Name:
		return a.binance.GetExchangeInfo(ctx, params)
	default:
		return exchanges.GetExchangeInfoResult{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Name)
	}
}
