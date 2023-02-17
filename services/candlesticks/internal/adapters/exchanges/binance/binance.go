package binance

import (
	"context"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/exchanges"
)

const Name = "binance"

type Service struct {
	client *client.Client
}

func New(c Config) (*Service, error) {
	return &Service{
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, c.Validate()
}

func (s *Service) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	service := s.client.NewKlinesService()

	// Set symbol
	service.Symbol(BinanceSymbol(payload.PairSymbol))

	// Set interval
	binanceInterval, err := PeriodToInterval(payload.Period)
	if err != nil {
		return nil, err
	}
	service.Interval(binanceInterval)

	// Set start and end time
	start, end := TimeToKLineTime(payload.Start), TimeToKLineTime(payload.End)
	service.StartTime(start)
	service.EndTime(end)

	// Set limit
	service.Limit(payload.Limit)

	// Get KLines
	kl, err := service.Do(ctx)
	if err != nil {
		return nil, err
	}

	// Change them to right format
	return KLinesToCandlesticks(payload.PairSymbol, payload.Period, kl, time.Now())
}
