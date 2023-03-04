//go:generate asyncapi-codegen -g application -p generated -i ../../../../api/asyncapi-spec/ticks.yaml -o ./generated/app.gen.go
//go:generate asyncapi-codegen -g client      -p generated -i ../../../../api/asyncapi-spec/ticks.yaml -o ./generated/client.gen.go
//go:generate asyncapi-codegen -g broker      -p generated -i ../../../../api/asyncapi-spec/ticks.yaml -o ./generated/broker.gen.go
//go:generate asyncapi-codegen -g types       -p generated -i ../../../../api/asyncapi-spec/ticks.yaml -o ./generated/types.gen.go
//go:generate asyncapi-codegen -g nats        -p generated -i ../../../../api/asyncapi-spec/ticks.yaml -o ./generated/nats.gen.go

package nats

import (
	"github.com/digital-feather/cryptellation/internal/ticks/app"
	"github.com/digital-feather/cryptellation/internal/ticks/ctrl/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc         *nats.Conn
	controller *generated.AppController
	ticks      app.Controller
}

func NewServer(c config.NATS, ticks app.Controller) (*Server, error) {
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
		nc:    nc,
		ticks: ticks,
	}, nil
}

func (s *Server) Listen() error {
	var err error

	s.controller, err = generated.NewAppController(generated.NewNATSController(s.nc))
	if err != nil {
		return err
	}

	return s.controller.SubscribeAll(newSubscriber(s.controller, s.ticks))
}

func (s *Server) Close() {
	if s.controller != nil {
		s.controller.Close()
	}
}
