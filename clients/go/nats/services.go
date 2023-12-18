package nats

import (
	"context"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/config"
)

type Services struct {
	Backtests    client.Backtests
	Candlesticks client.Candlesticks
	Exchanges    client.Exchanges
	Indicators   client.Indicators
	Ticks        client.Ticks
}

func NewServices(c config.NATS) (Services, error) {
	backtests, err := NewBacktests(c)
	if err != nil {
		return Services{}, err
	}

	candlesticks, err := NewCandlesticks(c)
	if err != nil {
		return Services{}, err
	}

	exchanges, err := NewExchanges(c)
	if err != nil {
		return Services{}, err
	}

	indicators, err := NewIndicators(c)
	if err != nil {
		return Services{}, err
	}

	ticks, err := NewTicks(c)
	if err != nil {
		return Services{}, err
	}

	return Services{
		Backtests:    backtests,
		Candlesticks: candlesticks,
		Exchanges:    exchanges,
		Indicators:   indicators,
		Ticks:        ticks,
	}, nil
}

func (s Services) Close(ctx context.Context) {
	s.Backtests.Close(ctx)
	s.Candlesticks.Close(ctx)
	s.Exchanges.Close(ctx)
	s.Indicators.Close(ctx)
	s.Ticks.Close(ctx)
}
