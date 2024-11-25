package binance

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/internal/activities"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/internal/domains/candlesticks/activities/exchanges/live/binance/entities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"go.temporal.io/sdk/worker"
)

// Activities is the struct that contains all the activities for the Binance exchange.
type Activities struct {
	*activities.Binance
}

// New creates a new Binance activities struct.
func New(c config.Binance) (*Activities, error) {
	s, err := activities.NewBinance(c)
	return &Activities{
		Binance: s,
	}, err
}

// Register registers the Binance activities with the given worker.
func (a Activities) Register(_ worker.Worker) {
	// No need to register the Binance activities, they are already registered
	// with its parent.
}

// GetCandlesticks gets the candlesticks for the given pair and period.
func (a *Activities) GetCandlesticks(
	ctx context.Context,
	params exchanges.GetCandlesticksParams,
) (exchanges.GetCandlesticksResult, error) {
	a.Client.Debug = true

	service := a.Client.NewKlinesService()

	// Set symbol
	service.Symbol(entities.BinanceSymbol(params.Pair))

	// Set interval
	binanceInterval, err := entities.PeriodToInterval(params.Period)
	if err != nil {
		return exchanges.GetCandlesticksResult{}, entities.WrapError(err)
	}
	service.Interval(binanceInterval)

	// Set start and end time
	service.StartTime(entities.TimeToKLineTime(params.Start))
	service.EndTime(entities.TimeToKLineTime(params.End))

	// Set limit
	if params.Limit > 0 {
		service.Limit(params.Limit)
	}

	// Get KLines
	kl, err := service.Do(ctx)
	if err != nil {
		return exchanges.GetCandlesticksResult{}, entities.WrapError(err)
	}

	// Change them to right format
	list, err := entities.KLinesToCandlesticks(params.Pair, params.Period, kl, time.Now())
	if err != nil {
		return exchanges.GetCandlesticksResult{}, entities.WrapError(err)
	}

	return exchanges.GetCandlesticksResult{
		List: list,
	}, nil
}
