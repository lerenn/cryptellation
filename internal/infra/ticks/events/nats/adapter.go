package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/internal/ctrl/ticks/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Adapter struct {
	broker extensions.BrokerController
	app    *events.AppController
	client client.Ticks
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
	client, err := natsClient.NewTicks(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		broker: broker,
		app:    app,
		client: client,
	}, nil
}

func (a *Adapter) Publish(ctx context.Context, tick tick.Tick) error {
	// Generated message
	msg := events.NewTickMessage()
	msg.Payload.Exchange = events.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = events.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = events.DateSchema(tick.Time)

	// Send message
	return a.app.PublishCryptellationTicksListenExchangePair(ctx,
		events.CryptellationTicksListenExchangePairParameters{
			Exchange: events.ExchangeNameSchema(tick.Exchange),
			Pair:     events.PairSymbolSchema(tick.PairSymbol),
		}, msg)
}

func (a *Adapter) Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(ctx, client.TicksFilterPayload{
		ExchangeName: "*",
		PairSymbol:   symbol,
	})
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
