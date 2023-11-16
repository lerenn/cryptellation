package ticks

import (
	"context"

	client "github.com/lerenn/cryptellation/clients/go"
	asyncapi "github.com/lerenn/cryptellation/pkg/asyncapi/ticks"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

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
