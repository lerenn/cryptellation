package domain

import (
	"context"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (t Ticks) Listen(ctx context.Context, exchange, pair string) (<-chan tick.Tick, error) {
	return t.events.Subscribe(ctx, pair)
}
