//go:generate asyncapi-codegen -g application -p generated -i ../../../../../api/asyncapi-spec/ticks.yaml -o ./generated/app.gen.go
//go:generate asyncapi-codegen -g client      -p generated -i ../../../../../api/asyncapi-spec/ticks.yaml -o ./generated/client.gen.go
//go:generate asyncapi-codegen -g broker      -p generated -i ../../../../../api/asyncapi-spec/ticks.yaml -o ./generated/broker.gen.go
//go:generate asyncapi-codegen -g types       -p generated -i ../../../../../api/asyncapi-spec/ticks.yaml -o ./generated/types.gen.go
//go:generate asyncapi-codegen -g nats        -p generated -i ../../../../../api/asyncapi-spec/ticks.yaml -o ./generated/nats.gen.go

package nats

import (
	"context"

	client "github.com/digital-feather/cryptellation/clients/go"
	natsClient "github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/digital-feather/cryptellation/internal/ticks/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/tick"
	"github.com/nats-io/nats.go"
)

type Adapter struct {
	nc     *nats.Conn
	app    *generated.AppController
	client client.Ticks
}

func New(c config.NATS) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Connect to NATS
	nc, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	// Create new app controller
	app, err := generated.NewAppController(generated.NewNATSController(nc))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client, err := natsClient.NewTicks(c)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		nc:     nc,
		app:    app,
		client: client,
	}, nil
}

func (a *Adapter) Publish(tick tick.Tick) error {
	// Generated message
	msg := generated.NewTickMessage()
	msg.Payload.Exchange = generated.ExchangeNameSchema(tick.Exchange)
	msg.Payload.PairSymbol = generated.PairSymbolSchema(tick.PairSymbol)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = generated.DateSchema(tick.Time)

	// Send message
	return a.app.PublishTicksListenExchangePair(generated.TicksListenExchangePairParameters{
		Exchange: generated.ExchangeNameSchema(tick.Exchange),
		Pair:     generated.PairSymbolSchema(tick.PairSymbol),
	}, msg)
}

func (a *Adapter) Subscribe(symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(context.Background(), client.TicksFilterPayload{
		ExchangeName: "*",
		PairSymbol:   symbol,
	})
}

func (a *Adapter) Close() {
	if a.app != nil {
		a.app.Close()
	}
}
