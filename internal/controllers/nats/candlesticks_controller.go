package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/candlesticks"
	"github.com/lerenn/cryptellation/internal/components/candlesticks"
	"github.com/lerenn/cryptellation/pkg/config"
)

type CandlesticksController struct {
	broker       *nats.Controller
	controller   *asyncapi.AppController
	logger       extensions.Logger
	candlesticks candlesticks.Interface
}

func NewCandlesticksController(c config.NATS, candlesticks candlesticks.Interface) (*CandlesticksController, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create a logger
	logger := loggers.NewECS()

	// Create an App controller
	controller, err := asyncapi.NewAppController(broker, asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &CandlesticksController{
		broker:       broker,
		controller:   controller,
		logger:       logger,
		candlesticks: candlesticks,
	}, nil
}

func (s *CandlesticksController) Listen() error {
	sub := newCandlesticksSubscriber(s.controller, s.candlesticks)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *CandlesticksController) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
