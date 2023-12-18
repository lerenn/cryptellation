// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=ticks.go -destination=mock/ticks.gen.go -package mock

package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/tick"
)

type Ticks interface {
	Register(ctx context.Context, payload TicksFilterPayload) error
	Listen(ctx context.Context, payload TicksFilterPayload) (<-chan tick.Tick, error)
	Unregister(ctx context.Context, payload TicksFilterPayload) error

	ServiceInfo(ctx context.Context) (ServiceInfo, error)
	Close(ctx context.Context)
}

type TicksFilterPayload struct {
	ExchangeName string
	PairSymbol   string
}
