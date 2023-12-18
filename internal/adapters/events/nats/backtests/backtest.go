package backtests

import (
	"context"

	asyncapi "github.com/lerenn/cryptellation/api/asyncapi/backtests"
	"github.com/lerenn/cryptellation/pkg/models/event"
)

func (a *Adapter) Publish(ctx context.Context, backtestID uint, evt event.Event) error {
	// Generated message
	msg := asyncapi.NewBacktestsEventMessage()

	// Set from event
	if err := msg.Set(evt); err != nil {
		return err
	}

	// Send message
	return a.app.PublishBacktestEvent(ctx, asyncapi.CryptellationBacktestsEventsParameters{
		Id: int64(backtestID),
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
