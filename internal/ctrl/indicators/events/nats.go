// SMA
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.16.0 -g application -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.16.0 -g client      -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.16.0 -g broker      -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.16.0 -g types       -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.16.0 -g nats        -p events -i ./../../../../api/asyncapi/indicators.yaml -o ./nats.gen.go

package events

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
	"github.com/lerenn/cryptellation/internal/core/indicators"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	nc         *nats.Conn
	controller *AppController
	indicators indicators.Interface
}

func NewNATS(c config.NATS, indicators indicators.Interface) (*NATS, error) {
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
		nc:         nc,
		indicators: indicators,
	}, nil
}

func (s *NATS) Listen() error {
	var err error

	s.controller, err = NewAppController(NewNATSController(s.nc))
	if err != nil {
		return err
	}
	s.controller.SetLogger(log.NewECS())

	return s.controller.SubscribeAll(context.Background(), newSubscriber(s.controller, s.indicators))
}

func (s *NATS) Close() {
	if s.controller != nil {
		s.controller.Close(context.Background())
	}
}
