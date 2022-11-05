package application

import (
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/backtests"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

type Application struct {
	Backtests backtests.Operator
}

func New(cs candlesticks.Client) (*Application, error) {
	repository, err := redis.New()
	if err != nil {
		return nil, err
	}

	ps, err := nats.New()
	if err != nil {
		return nil, err
	}

	return &Application{
		Backtests: backtests.New(repository, ps, cs),
	}, nil
}
