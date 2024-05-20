// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=ticks.go -destination=mock.gen.go -package client

package client

import (
	"context"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/event"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Client interface {
	SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error)
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type TicksFilterPayload struct {
	Exchange string
	Pair     string
}
