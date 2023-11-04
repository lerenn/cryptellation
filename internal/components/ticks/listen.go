package ticks

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/tick"
)

func (t Ticks) Listen(ctx context.Context, exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.events.Subscribe(ctx, pairSymbol)
}
