package backtests

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/event"

	asyncapi "github.com/lerenn/cryptellation/svc/backtests/api/asyncapi"

	"github.com/google/uuid"
)

func (a *Adapter) Publish(ctx context.Context, backtestID uuid.UUID, evt event.Event) error {
	// Generated message
	msg := asyncapi.NewEventMessage()

	// Set from event
	if err := msg.Set(evt); err != nil {
		return err
	}

	// Send message
	return a.app.SendAsEventOperation(ctx, asyncapi.EventsChannelParameters{
		Id: backtestID.String(),
	}, msg)
}

func (a *Adapter) Subscribe(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error) {
	return a.client.ListenEvents(ctx, backtestID)
}

func (a *Adapter) Close(ctx context.Context) {
	if a.app != nil {
		a.app.Close(ctx)
	}
}
