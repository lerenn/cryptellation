package app

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/event"
)

type Ticks interface {
	ListeningNotificationReceived(ctx context.Context, ts event.TickSubscription)
	Close(ctx context.Context)
}
