package binance

import (
	"context"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/internal/core/candlesticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

const Name = "binance"

type Service struct {
	client *client.Client
}

func New(c config.Binance) (*Service, error) {
	return &Service{
		client: client.NewClient(
			c.ApiKey,
			c.SecretKey),
	}, c.Validate()
}

func (s *Service) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	s.client.Debug = true

	service := s.client.NewKlinesService()

	// Set symbol
	service.Symbol(BinanceSymbol(payload.PairSymbol))

	// Set interval
	binanceInterval, err := PeriodToInterval(payload.Period)
	if err != nil {
		return nil, wrapError(err)
	}
	service.Interval(binanceInterval)

	// Set start and end time
	service.StartTime(TimeToKLineTime(payload.Start))
	service.EndTime(TimeToKLineTime(payload.End))

	// Set limit
	if payload.Limit > 0 {
		service.Limit(payload.Limit)
	}

	// Get KLines
	kl, err := service.Do(ctx)
	if err != nil {
		return nil, wrapError(err)
	}

	// Change them to right format
	return KLinesToCandlesticks(payload.PairSymbol, payload.Period, kl, time.Now())
}
