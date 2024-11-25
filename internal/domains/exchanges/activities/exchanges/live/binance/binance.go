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

type Activities struct {
	*activities.Binance
}

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

func (a *Activities) GetExchangeInfo(ctx context.Context, params exchanges.GetExchangeInfoParams) (exchanges.GetExchangeInfoResult, error) {
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
