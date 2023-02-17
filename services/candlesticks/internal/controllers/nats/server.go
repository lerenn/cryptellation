package nats

import (
	async "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/async"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc           *nats.Conn
	controller   *async.AppController
	candlesticks candlesticks.Port
}

func NewServer(c Config, candlesticks candlesticks.Port) (*Server, error) {
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

	s.controller, err = async.NewAppController(async.NewNATSController(s.nc))
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
