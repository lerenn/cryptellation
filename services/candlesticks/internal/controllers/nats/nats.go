//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g application -p gen -i ../../../api/asyncapi.yaml -o ./gen/app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g broker      -p gen -i ../../../api/asyncapi.yaml -o ./gen/broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g types       -p gen -i ../../../api/asyncapi.yaml -o ./gen/types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.5.0 -g nats        -p gen -i ../../../api/asyncapi.yaml -o ./gen/nats.gen.go

package nats

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/nats/gen"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	controller *gen.AppController
	app        *application.Application
}

func New(app *application.Application) *NATS {
	return &NATS{
		app: app,
	}
}

func (n *NATS) Run() error {
	// Load config
	config := loadConfig()
	if err := config.Validate(); err != nil {
		return err
	}

	// Connect to NATS
	nc, err := nats.Connect(config.URL())
	if err != nil {
		return err
	}

	// Create a new application controller
	n.controller, err = gen.NewAppController(gen.NewNATSController(nc))
	if err != nil {
		return err
	}

	return n.controller.SubscribeAll(newSubscriber(n.controller, n.app))
}

func (n *NATS) Close() {
	if n.controller != nil {
		n.controller.Close()
	}
}
