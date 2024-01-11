package nats

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"
	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"
	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
)

type Services struct {
	Backtests    backtests.Client
	Candlesticks candlesticks.Client
	Exchanges    exchanges.Client
	Indicators   indicators.Client
	Ticks        ticks.Client
}

func NewServices(c config.NATS) (Services, error) {
	backtests, err := backtests.NewClient(c)
	if err != nil {
		return Services{}, err
	}

	candlesticks, err := candlesticks.NewClient(c)
	if err != nil {
		return Services{}, err
	}

	exchanges, err := exchanges.NewClient(c)
	if err != nil {
		return Services{}, err
	}

	indicators, err := indicators.NewClient(c)
	if err != nil {
		return Services{}, err
	}

	ticks, err := ticks.NewClient(c)
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
