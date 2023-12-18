package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/exchanges"
	"github.com/lerenn/cryptellation/internal/components/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
)

type ExchangesController struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	exchanges  exchanges.Interface
}

func NewExchangesController(c config.NATS, exchanges exchanges.Interface) (*ExchangesController, error) {
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

	return &ExchangesController{
		broker:     broker,
		controller: controller,
		exchanges:  exchanges,
		logger:     logger,
	}, nil
}

func (s *ExchangesController) Listen() error {
	sub := newExchangesSubscriber(s.controller, s.exchanges)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *ExchangesController) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
