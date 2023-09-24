// SMA
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g application -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g user        -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g types       -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./types.gen.go

package events

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/cryptellation/internal/core/indicators"
	"github.com/lerenn/cryptellation/pkg/config"
)

type NATS struct {
	broker     *nats.Controller
	controller *AppController
	logger     extensions.Logger
	indicators indicators.Interface
}

func NewNATS(c config.NATS, indicators indicators.Interface) (*NATS, error) {
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
		indicators: indicators,
		logger:     logger,
	}, nil
}

func (s *NATS) Listen() error {
	sub := newSubscriber(s.controller, s.indicators)
	return s.controller.SubscribeAll(context.Background(), sub)
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
