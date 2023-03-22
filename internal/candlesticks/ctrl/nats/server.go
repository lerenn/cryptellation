package nats

import (
	"github.com/digital-feather/cryptellation/internal/candlesticks/app"
	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc           *nats.Conn
	controller   *generated.AppController
	candlesticks app.Controller
}

func NewServer(c config.NATS, candlesticks app.Controller) (*Server, error) {
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
		nc:           nc,
		candlesticks: candlesticks,
	}, nil
}

func (s *Server) Listen() error {
	var err error

	s.controller, err = generated.NewAppController(generated.NewNATSController(s.nc))
	if err != nil {
		return err
	}

	return s.controller.SubscribeAll(newSubscriber(s.controller, s.candlesticks))
}

func (s *Server) Close() {
	if s.controller != nil {
		s.controller.Close()
	}
}
