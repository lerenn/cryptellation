package live

import (
	"context"
	"fmt"

	binancePkg "github.com/lerenn/cryptellation/v1/internal/activities"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges/live/binance"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

// Activities is the live exchanges activities, regrouping all exchanges activities.
type Activities struct {
	binance exchanges.Interface
}

// New will create a new live activities.
func New() (Activities, error) {
	b, err := binance.New(config.LoadBinanceTest())
	if err != nil {
		return Activities{}, err
	}

	return Activities{
		binance: b,
	}, nil
}

// Register will register the activities.
func (a Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.GetCandlesticks, activity.RegisterOptions{Name: exchanges.GetCandlesticksActivityName})
}

// GetCandlesticks will get the candlesticks.
func (e Activities) GetCandlesticks(
	ctx context.Context,
	params exchanges.GetCandlesticksParams,
) (exchanges.GetCandlesticksResult, error) {
	switch params.Exchange {
	case binancePkg.BinanceInfos.Name:
		res, err := e.binance.GetCandlesticks(ctx, params)
		return res, err
	default:
		return exchanges.GetCandlesticksResult{}, fmt.Errorf("%w: %q", exchanges.ErrInexistantExchange, params.Exchange)
	}
}
