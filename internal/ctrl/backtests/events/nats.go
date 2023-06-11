package events

import (
	"github.com/lerenn/cryptellation/internal/core/backtests"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	nc         *nats.Conn
	controller *AppController
	exchanges  backtests.Interface
}

func NewNATS(c config.NATS, exchanges backtests.Interface) (*NATS, error) {
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
		nc:        nc,
		exchanges: exchanges,
	}, nil
}

func (s *NATS) Listen() error {
	var err error

	s.controller, err = NewAppController(NewNATSController(s.nc))
	if err != nil {
		return err
	}

	return s.controller.SubscribeAll(newSubscriber(s.controller, s.exchanges))
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close()
	}
}
