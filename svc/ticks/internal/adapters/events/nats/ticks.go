package nats

import (
	"context"

	asyncapi "github.com/lerenn/cryptellation/svc/ticks/api/asyncapi"
	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (a *Adapter) Publish(ctx context.Context, tick tick.Tick) error {
	// Generated message
	msg := asyncapi.NewTickMessage()
	msg.Payload.Exchange = asyncapi.ExchangeSchema(tick.Exchange)
	msg.Payload.Pair = asyncapi.PairSchema(tick.Pair)
	msg.Payload.Price = tick.Price
	msg.Payload.Time = asyncapi.DateSchema(tick.Time)

	// Send message
	return a.app.SendAsLiveOperation(ctx,
		asyncapi.LiveChannelParameters{
			Exchange: tick.Exchange,
			Pair:     tick.Pair,
		}, msg)
}

func (a *Adapter) Subscribe(ctx context.Context, symbol string) (<-chan tick.Tick, error) {
	return a.client.Listen(ctx, client.TicksFilterPayload{
		Exchange: "*",
		Pair:     symbol,
	})
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
