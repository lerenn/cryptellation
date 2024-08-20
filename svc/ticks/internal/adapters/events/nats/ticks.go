package nats

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/event"

	asyncapi "github.com/lerenn/cryptellation/ticks/api/asyncapi"
	"github.com/lerenn/cryptellation/ticks/pkg/tick"
)

func (a *Adapter) PublishTick(ctx context.Context, tick tick.Tick) error {
	// Generated message
	msg := asyncapi.NewTickMessage()
	msg.FromModel(tick)

	// Send message
	return a.app.SendAsSendNewTicksOperation(ctx,
		asyncapi.LiveChannelParameters{
			Exchange: tick.Exchange,
			Pair:     tick.Pair,
		}, msg)
}

func (a *Adapter) SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick, 16)
	err := a.user.SubscribeToSendNewTicksOperation(ctx, asyncapi.LiveChannelParameters{
		Exchange: sub.Exchange,
		Pair:     sub.Pair,
	}, func(ctx context.Context, msg asyncapi.TickMessage) error {
		ch <- msg.ToModel()
		return nil
	})
	return ch, err
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
