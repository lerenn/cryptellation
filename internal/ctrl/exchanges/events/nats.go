package events

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/cryptellation/internal/core/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
)

type NATS struct {
	broker     *nats.Controller
	controller *AppController
	logger     extensions.Logger
	exchanges  exchanges.Interface
}

func NewNATS(c config.NATS, exchanges exchanges.Interface) (*NATS, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create an App controller
	controller, err := NewAppController(broker, WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &NATS{
		broker:     broker,
		controller: controller,
		exchanges:  exchanges,
		logger:     logger,
	}, nil
}

func (s *NATS) Listen() error {
	sub := newSubscriber(s.controller, s.exchanges)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
