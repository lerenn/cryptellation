package app

import (
	"context"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Ticks interface {
	Listen(ctx context.Context, exchange, pair string) (<-chan tick.Tick, error)
	Register(ctx context.Context, exchange, pair string) (int64, error)
	Unregister(ctx context.Context, exchange, pair string) (int64, error)
}
