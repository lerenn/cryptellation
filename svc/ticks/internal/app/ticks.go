package app

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Ticks interface {
	ListeningNotificationReceived(ctx context.Context, ts event.PricesSubscription)
	Close(ctx context.Context)
}
