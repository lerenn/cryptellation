package app

import (
	"context"

	"cryptellation/pkg/models/event"
)

type Ticks interface {
	ListeningNotificationReceived(ctx context.Context, ts event.TickSubscription)
	Close(ctx context.Context)
}
