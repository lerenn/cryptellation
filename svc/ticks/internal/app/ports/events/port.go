// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package events

package events

import (
	"context"

	"cryptellation/pkg/models/event"

	"cryptellation/svc/ticks/pkg/tick"
)

type Port interface {
	PublishTick(ctx context.Context, tick tick.Tick) error
	SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error)

	Close(ctx context.Context)
}
