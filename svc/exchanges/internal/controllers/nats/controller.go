package nats

import (
	"context"

	asyncapipkg "cryptellation/internal/asyncapi"
	"cryptellation/pkg/config"

	asyncapi "cryptellation/svc/exchanges/api/asyncapi"
	exchanges "cryptellation/svc/exchanges/internal/app"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type Controller struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	exchanges  exchanges.Exchanges
}

func NewController(c config.NATS, exchanges exchanges.Exchanges) (*Controller, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker, err := nats.NewController(c.URL())
	if err != nil {
		return nil, err
	}

	// Create an App controller
	logger := asyncapipkg.LoggerWrapper{}
	controller, err := asyncapi.NewAppController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &Controller{
		broker:     broker,
		controller: controller,
		exchanges:  exchanges,
		logger:     logger,
	}, nil
}

func (s *Controller) Listen() error {
	sub := newSubscriber(s.controller, s.exchanges)
	return s.controller.SubscribeToAllChannels(context.Background(), sub)
}

func (s *Controller) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
