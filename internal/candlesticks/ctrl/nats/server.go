//go:generate asyncapi-codegen -g application -p generated -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./generated/app.gen.go
//go:generate asyncapi-codegen -g client      -p generated -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./generated/client.gen.go
//go:generate asyncapi-codegen -g broker      -p generated -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./generated/broker.gen.go
//go:generate asyncapi-codegen -g types       -p generated -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./generated/types.gen.go
//go:generate asyncapi-codegen -g nats        -p generated -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./generated/nats.gen.go

package nats

import (
	"github.com/digital-feather/cryptellation/internal/candlesticks/app"
	"github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats/generated"
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
