//go:generate asyncapi-codegen -g application -p internal -i ../../../../api/asyncapi-spec/exchanges.yaml -o ./internal/app.gen.go
//go:generate asyncapi-codegen -g client      -p internal -i ../../../../api/asyncapi-spec/exchanges.yaml -o ./internal/client.gen.go
//go:generate asyncapi-codegen -g broker      -p internal -i ../../../../api/asyncapi-spec/exchanges.yaml -o ./internal/broker.gen.go
//go:generate asyncapi-codegen -g types       -p internal -i ../../../../api/asyncapi-spec/exchanges.yaml -o ./internal/types.gen.go
//go:generate asyncapi-codegen -g nats        -p internal -i ../../../../api/asyncapi-spec/exchanges.yaml -o ./internal/nats.gen.go

package nats

import (
	"github.com/digital-feather/cryptellation/internal/exchanges/app"
	"github.com/digital-feather/cryptellation/internal/exchanges/ctrl/nats/internal"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc         *nats.Conn
	controller *internal.AppController
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

	s.controller, err = internal.NewAppController(internal.NewNATSController(s.nc))
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
