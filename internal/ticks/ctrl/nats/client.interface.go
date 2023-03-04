// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=client.interface.go -destination=client.mock.gen.go -package nats

package nats

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/tick"
)

type Client interface {
	Register(ctx context.Context, payload TicksFilterPayload) error
	Listen(ctx context.Context, payload TicksFilterPayload) (<-chan tick.Tick, error)
	Unregister(ctx context.Context, payload TicksFilterPayload) error
	Close()
}

type TicksFilterPayload struct {
	ExchangeName string
	PairSymbol   string
}
