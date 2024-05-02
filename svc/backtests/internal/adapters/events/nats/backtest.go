package backtests

import (
	"context"
	"fmt"

	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
)

func (a *Adapter) Publish(ctx context.Context, backtestID uint, evt event.Event) error {
	// Generated message
	msg := asyncapi.NewEventMessage()

	// Set from event
	if err := msg.Set(evt); err != nil {
		return err
	}

	// Send message
	return a.app.SendAsEventOperation(ctx, asyncapi.EventsChannelParameters{
		Id: fmt.Sprintf("%d", backtestID),
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
