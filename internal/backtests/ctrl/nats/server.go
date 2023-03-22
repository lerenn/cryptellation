package nats

import (
	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/backtests"
	"github.com/digital-feather/cryptellation/internal/backtests/app"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc         *nats.Conn
	controller *asyncapi.AppController
	exchanges  app.Controller
}

func NewServer(c config.NATS, exchanges app.Controller) (*Server, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Connect to NATS
	nc, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	return &Server{
		nc:        nc,
		exchanges: exchanges,
	}, nil
}

func (s *Server) Listen() error {
	var err error

	s.controller, err = asyncapi.NewAppController(asyncapi.NewNATSController(s.nc))
	if err != nil {
		return err
	}

	return s.controller.SubscribeAll(newSubscriber(s.controller, s.exchanges))
}

func (s *Server) Close() {
	if s.controller != nil {
		s.controller.Close()
	}
}
