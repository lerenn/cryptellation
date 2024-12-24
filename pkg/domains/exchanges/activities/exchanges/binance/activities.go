package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
	"go.temporal.io/sdk/worker"
)

// Activities is the live Binance activities.
type Activities struct {
	*activities.Binance
}

// New creates a new Binance activities.
func New() (*Activities, error) {
	s, err := activities.NewBinance(config.LoadBinanceTest())
	return &Activities{
		Binance: s,
	}, err
}

// Name returns the name of the Binance activities.
func (a *Activities) Name() string {
	return activities.BinanceInfos.Name
}

// Register registers the Binance activities with the given worker.
func (a *Activities) Register(_ worker.Worker) {
	// No need to register the Binance activities, they are already registered
	// with its parent.
}

// ListExchangesActivity returns the names of the exchanges.
func (a *Activities) ListExchangesActivity(
	_ context.Context,
	_ exchanges.ListExchangesActivityParams,
) (exchanges.ListExchangesActivityResults, error) {
	return exchanges.ListExchangesActivityResults{
		List: []string{
			activities.BinanceInfos.Name,
		},
	}, nil
}

// GetExchangeActivity returns the exchange information for the given exchange.
func (a *Activities) GetExchangeActivity(
	ctx context.Context,
	_ exchanges.GetExchangeActivityParams,
) (exchanges.GetExchangeActivityResults, error) {
	exchangeInfos, err := a.Client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchanges.GetExchangeActivityResults{}, err
	}

	pairs := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairs[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := activities.BinanceInfos
	exch.Pairs = pairs
	exch.LastSyncTime = time.Now()

	return exchanges.GetExchangeActivityResults{
		Exchange: exch,
	}, nil
}
