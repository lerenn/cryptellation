package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/candlesticks/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
)

type Controller struct {
	broker       *nats.Controller
	controller   *asyncapi.AppController
	logger       extensions.Logger
	candlesticks app.Candlesticks
}

func NewController(c config.NATS, candlesticks app.Candlesticks) (*Controller, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker, err := nats.NewController(c.URL())
	if err != nil {
		return nil, err
	}

	// Create a logger
	logger := loggers.NewECS()

	// Create an App controller
	controller, err := asyncapi.NewAppController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &Controller{
		broker:       broker,
		controller:   controller,
		logger:       logger,
		candlesticks: candlesticks,
	}, nil
}

func (s *Controller) Listen() error {
	sub := newCandlesticksSubscriber(s.controller, s.candlesticks)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *Controller) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
