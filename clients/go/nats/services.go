package nats

import (
	"context"
	"fmt"

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
	// Check the configuration before creating clients
	if err := c.Validate(); err != nil {
		return Services{}, err
	}

	backtests, err := backtests.NewClient(c)
	if err != nil {
		return Services{}, fmt.Errorf("error when creating new backtests client: %w", err)
	}

	candlesticks, err := candlesticks.NewClient(c)
	if err != nil {
		return Services{}, fmt.Errorf("error when creating new candlesticks client: %w", err)
	}

	exchanges, err := exchanges.NewClient(c)
	if err != nil {
		return Services{}, fmt.Errorf("error when creating new exchanges client: %w", err)
	}

	indicators, err := indicators.NewClient(c)
	if err != nil {
		return Services{}, fmt.Errorf("error when creating new indicators client: %w", err)
	}

	ticks, err := ticks.NewClient(c)
	if err != nil {
		return Services{}, fmt.Errorf("error when creating new ticks client: %w", err)
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
