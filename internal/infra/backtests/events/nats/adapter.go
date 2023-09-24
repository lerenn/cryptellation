package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/internal/ctrl/backtests/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Adapter struct {
	broker extensions.BrokerController
	app    *events.AppController
	client client.Backtests
}

func New(c config.NATS) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Create a NATS Controller
	broker := nats.NewController(c.URL())

	// Create new app controller
	app, err := events.NewAppController(broker)
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewBacktests(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		broker: broker,
		app:    app,
		client: client,
	}, nil
}

func (a *Adapter) Publish(ctx context.Context, backtestID uint, evt event.Event) error {
	// Generated message
	msg := events.NewBacktestsEventMessage()

	// Set from event
	if err := msg.Set(evt); err != nil {
		return err
	}

	// Send message
	return a.app.PublishCryptellationBacktestsEventsID(ctx, events.CryptellationBacktestsEventsIDParameters{
		ID: int64(backtestID),
	}, msg)
}

func (a *Adapter) Subscribe(ctx context.Context, backtestID uint) (<-chan event.Event, error) {
	return a.client.ListenEvents(ctx, backtestID)
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
