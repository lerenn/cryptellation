package backtests

import (
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	adapter "github.com/lerenn/cryptellation/pkg/adapters/events/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	natsClient "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
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
	logger := loggers.NewECS()
	app, err := asyncapi.NewAppController(events.Broker(), asyncapi.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewClient(c, natsClient.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return &Adapter{
		events: events,
		app:    app,
		client: client,
	}, nil
}
