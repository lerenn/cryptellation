package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"github.com/lerenn/cryptellation/v1/internal/domains/exchanges/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/config"
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

// Register registers the Binance activities with the given worker.
func (a Activities) Register(_ worker.Worker) {
	// No need to register the Binance activities, they are already registered
	// with its parent.
}

// ListExchangesNames returns the names of the exchanges.
func (a *Activities) ListExchangesNames(
	_ context.Context,
	_ exchanges.ListExchangesNamesParams,
) (exchanges.ListExchangesNamesResult, error) {
	return exchanges.ListExchangesNamesResult{
		List: []string{
			activities.BinanceInfos.Name,
		},
	}, nil
}

// GetExchangeInfo returns the exchange information for the given exchange.
func (a *Activities) GetExchangeInfo(
	ctx context.Context,
	_ exchanges.GetExchangeInfoParams,
) (exchanges.GetExchangeInfoResult, error) {
	exchangeInfos, err := a.Client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchanges.GetExchangeInfoResult{}, err
	}

	pairs := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairs[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := activities.BinanceInfos
	exch.Pairs = pairs
	exch.LastSyncTime = time.Now()

	return exchanges.GetExchangeInfoResult{
		Exchange: exch,
	}, nil
}
