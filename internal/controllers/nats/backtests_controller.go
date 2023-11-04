package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/cryptellation/internal/components/backtests"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/backtests"
	"github.com/lerenn/cryptellation/pkg/config"
)

type BacktestsController struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	backtests  backtests.Interface
}

func NewBacktestsController(c config.NATS, backtests backtests.Interface) (*BacktestsController, error) {
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

	return &BacktestsController{
		broker:     broker,
		controller: controller,
		backtests:  backtests,
		logger:     logger,
	}, nil
}

func (s *BacktestsController) Listen() error {
	sub := newBacktestsSubscriber(s.controller, s.backtests)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *BacktestsController) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
