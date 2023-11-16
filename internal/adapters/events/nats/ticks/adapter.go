package ticks

import (
	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	adapter "github.com/lerenn/cryptellation/internal/adapters/events/nats"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/ticks"
	"github.com/lerenn/cryptellation/pkg/config"
)

type Adapter struct {
	events *adapter.Adapter
	app    *asyncapi.AppController
	client client.Ticks
}

func New(c config.NATS) (*Adapter, error) {
	// Create embedded database access
	events, err := adapter.New(c)
	if err != nil {
		return nil, err
	}

	// Create new app controller
	app, err := asyncapi.NewAppController(events.Broker())
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewTicks(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		events: events,
		app:    app,
		client: client,
	}, nil
}
