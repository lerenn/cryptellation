package ticks

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/tick"
)

type Interface interface {
	Listen(exchange, pairSymbol string) (<-chan tick.Tick, error)
	Register(ctx context.Context, exchange, pairSymbol string) (int64, error)
	Unregister(ctx context.Context, exchange, pairSymbol string) (int64, error)
}
