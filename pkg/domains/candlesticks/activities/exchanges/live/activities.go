package live

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges/live/binance"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the live exchanges activities, regrouping all exchanges activities.
type Activities struct {
	binance exchanges.Exchanges
}

// New will create a new live activities.
func New() (*Activities, error) {
	b, err := binance.New(config.LoadBinanceTest())
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
		a.GetCandlesticksActivity,
		activity.RegisterOptions{Name: exchanges.GetCandlesticksActivityName},
	)
}

// GetCandlesticksActivity will get the candlesticks.
func (a *Activities) GetCandlesticksActivity(
	ctx context.Context,
	params exchanges.GetCandlesticksActivityParams,
) (exchanges.GetCandlesticksActivityResults, error) {
	switch params.Exchange {
	case binancePkg.BinanceInfos.Name:
		res, err := a.binance.GetCandlesticksActivity(ctx, params)
		return res, err
	default:
		err := fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
		return exchanges.GetCandlesticksActivityResults{}, err
	}
}
