package nats

import (
	"context"

	asyncapipkg "cryptellation/internal/asyncapi"
	"cryptellation/internal/config"

	asyncapi "cryptellation/svc/indicators/api/asyncapi"
	"cryptellation/svc/indicators/internal/app"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type Controller struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	indicators app.Indicators
}

func NewController(c config.NATS, indicators app.Indicators) (*Controller, error) {
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
		indicators: indicators,
		logger:     logger,
	}, nil
}

func (s *Controller) Listen() error {
	sub := newSubscriber(s.controller, s.indicators)
	return s.controller.SubscribeToAllChannels(context.Background(), sub)
}

func (s *Controller) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
