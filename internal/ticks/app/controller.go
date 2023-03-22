package app

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/types/tick"
)

type Controller interface {
	Listen(exchange, pairSymbol string) (<-chan tick.Tick, error)
	Register(ctx context.Context, exchange, pairSymbol string) (int64, error)
	Unregister(ctx context.Context, exchange, pairSymbol string) (int64, error)
}
