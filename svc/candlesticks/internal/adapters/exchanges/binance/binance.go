package binance

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/internal/adapters/exchanges/binance"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/svc/candlesticks/internal/adapters/exchanges/binance/entities"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type Service struct {
	*binance.Service
}

func New(c config.Binance) (*Service, error) {
	s, err := binance.New(c)
	return &Service{
		Service: s,
	}, err
}

func (s *Service) GetCandlesticks(ctx context.Context, payload exchanges.GetCandlesticksPayload) (*candlestick.List, error) {
	s.Client.Debug = true

	service := s.Client.NewKlinesService()

	// Set symbol
	service.Symbol(entities.BinanceSymbol(payload.Pair))

	// Set interval
	binanceInterval, err := entities.PeriodToInterval(payload.Period)
	if err != nil {
		return nil, entities.WrapError(err)
	}
	service.Interval(binanceInterval)

	// Set start and end time
	service.StartTime(entities.TimeToKLineTime(payload.Start))
	service.EndTime(entities.TimeToKLineTime(payload.End))

	// Set limit
	if payload.Limit > 0 {
		service.Limit(payload.Limit)
	}

	// Get KLines
	kl, err := service.Do(ctx)
	if err != nil {
		return nil, entities.WrapError(err)
	}
	for _, k := range kl {
		telemetry.L(ctx).Debugf("Received Binance KLine: %+v", k)
	}

	// Change them to right format
	return entities.KLinesToCandlesticks(payload.Pair, payload.Period, kl, time.Now())
}
