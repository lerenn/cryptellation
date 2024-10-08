package nats

import (
	"context"

	asyncapipkg "github.com/lerenn/cryptellation/internal/asyncapi"

	"github.com/lerenn/cryptellation/pkg/config"

	asyncapi "github.com/lerenn/cryptellation/svc/forwardtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type Controller struct {
	broker       *nats.Controller
	controller   *asyncapi.AppController
	logger       extensions.Logger
	forwardtests app.ForwardTests
}

func NewController(c config.NATS, forwardtests app.ForwardTests) (*Controller, error) {
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
		broker:       broker,
		controller:   controller,
		forwardtests: forwardtests,
		logger:       logger,
	}, nil
}

func (s *Controller) Listen() error {
	sub := newSubscriber(s.controller, s.forwardtests)
	return s.controller.SubscribeToAllChannels(context.Background(), sub)
}

func (s *Controller) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
