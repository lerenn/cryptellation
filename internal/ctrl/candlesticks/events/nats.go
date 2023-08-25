package events

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
	"github.com/lerenn/cryptellation/internal/core/candlesticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	nc           *nats.Conn
	controller   *AppController
	candlesticks candlesticks.Interface
}

func NewNATS(c config.NATS, candlesticks candlesticks.Interface) (*NATS, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Connect to NATS
	nc, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	return &NATS{
		nc:           nc,
		candlesticks: candlesticks,
	}, nil
}

func (s *NATS) Listen() error {
	var err error

	s.controller, err = NewAppController(NewNATSController(s.nc))
	if err != nil {
		return err
	}
	s.controller.SetLogger(log.NewECS())

	return s.controller.SubscribeAll(context.Background(), newSubscriber(s.controller, s.candlesticks))
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
