//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g application -p internal -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./internal/app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g client      -p internal -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./internal/client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g broker      -p internal -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./internal/broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g types       -p internal -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./internal/types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g nats        -p internal -i ../../../../api/asyncapi-spec/candlesticks.yaml -o ./internal/nats.gen.go

package nats

import (
	"github.com/digital-feather/cryptellation/internal/components/candlesticks"
	async "github.com/digital-feather/cryptellation/internal/controllers/nats"
	"github.com/digital-feather/cryptellation/internal/controllers/nats/candlesticks/internal"
	"github.com/nats-io/nats.go"
)

type Server struct {
	nc           *nats.Conn
	controller   *internal.AppController
	candlesticks candlesticks.Port
}

func NewServer(c async.Config, candlesticks candlesticks.Port) (*Server, error) {
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

	s.controller, err = internal.NewAppController(internal.NewNATSController(s.nc))
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
