package nats

import (
	adapter "github.com/lerenn/cryptellation/pkg/adapters/events/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/ticks/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	natsClient "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
)

type Adapter struct {
	events *adapter.Adapter
	app    *asyncapi.AppController
	client client.Client
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
	client, err := natsClient.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		events: events,
		app:    app,
		client: client,
	}, nil
}
