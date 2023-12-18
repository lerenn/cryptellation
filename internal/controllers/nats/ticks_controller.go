package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/ticks"
	"github.com/lerenn/cryptellation/internal/components/ticks"
	"github.com/lerenn/cryptellation/pkg/config"
)

type TicksController struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	ticks      ticks.Interface
}

func NewTicksController(c config.NATS, ticks ticks.Interface) (*TicksController, error) {
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

	return &TicksController{
		broker:     broker,
		controller: controller,
		ticks:      ticks,
		logger:     logger,
	}, nil
}

func (s *TicksController) Listen() error {
	sub := newTicksSubscriber(s.controller, s.ticks)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *TicksController) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
