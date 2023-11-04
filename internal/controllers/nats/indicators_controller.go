package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/cryptellation/internal/components/indicators"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/indicators"
	"github.com/lerenn/cryptellation/pkg/config"
)

type IndicatorsSubscriber struct {
	broker     *nats.Controller
	controller *asyncapi.AppController
	logger     extensions.Logger
	indicators indicators.Interface
}

func NewIndicatorsSubscriber(c config.NATS, indicators indicators.Interface) (*IndicatorsSubscriber, error) {
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

	return &IndicatorsSubscriber{
		broker:     broker,
		controller: controller,
		indicators: indicators,
		logger:     logger,
	}, nil
}

func (s *IndicatorsSubscriber) Listen() error {
	sub := newIndicatorsSubscriber(s.controller, s.indicators)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *IndicatorsSubscriber) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
