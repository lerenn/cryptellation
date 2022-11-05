package application

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ticks"
)

type Application struct {
	Ticks ticks.Operator
}

func New(exchanges map[string]exchanges.Port) (*Application, error) {
	repository, err := redis.New()
	if err != nil {
		return nil, err
	}

	if err := repository.ClearSymbolListenersCount(context.TODO()); err != nil {
		return nil, err
	}

	ps, err := nats.New()
	if err != nil {
		return nil, err
	}

	return &Application{
		Ticks: ticks.New(ps, repository, exchanges),
	}, nil
}
