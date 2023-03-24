package ticks

import (
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/services/ticks"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	nc         *nats.Conn
	controller *AppController
	ticks      ticks.Interface
}

func NewNATS(c config.NATS, ticks ticks.Interface) (*NATS, error) {
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
		nc:    nc,
		ticks: ticks,
	}, nil
}

func (s *NATS) Listen() error {
	var err error

	s.controller, err = NewAppController(NewNATSController(s.nc))
	if err != nil {
		return err
	}

	return s.controller.SubscribeAll(newSubscriber(s.controller, s.ticks))
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close()
	}
}
