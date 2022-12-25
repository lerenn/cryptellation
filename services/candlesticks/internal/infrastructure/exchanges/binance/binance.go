package binance

import (
	"context"
	"fmt"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

const Name = "binance"

type Service struct {
	config Config
	client *client.Client
}

func New() (*Service, error) {
	var c Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading binance config: %w", err)
	}

	return &Service{
		config: c,
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, nil
}

// func (s *Service) Candlesticks(pairSymbol string, per period.Symbol) (exchanges.CandlesticksService, error) {
// 	service := s.client.NewKlinesService()
// 	service.Symbol(BinanceSymbol(pairSymbol))

// 	binanceInterval, err := PeriodToInterval(per)
// 	if err != nil {
// 		return nil, err
// 	}
// 	service.Interval(binanceInterval)

// 	return &CandlestickService{
// 		service:    service,
// 		pairSymbol: pairSymbol,
// 		period:     per,
// 	}, nil
// }

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
