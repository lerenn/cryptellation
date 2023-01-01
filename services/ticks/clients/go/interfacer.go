// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=interfacer.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
)

type Interfacer interface {
	Register(ctx context.Context, exchange, symbol string) error
	Unregister(ctx context.Context, exchange, symbol string) error
	Listen(symbol string) (<-chan tick.Tick, error)
}
