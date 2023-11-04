package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/ticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Adapter struct {
	broker extensions.BrokerController
	app    *asyncapi.AppController
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
	app, err := asyncapi.NewAppController(broker)
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
	msg := asyncapi.NewTickMessage()
	msg.Payload.Exchange = asyncapi.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = asyncapi.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = asyncapi.DateSchema(tick.Time)

	// Send message
	return a.app.PublishCryptellationTicksListenExchangePair(ctx,
		asyncapi.CryptellationTicksListenExchangePairParameters{
			Exchange: asyncapi.ExchangeNameSchema(tick.Exchange),
			Pair:     asyncapi.PairSymbolSchema(tick.PairSymbol),
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
