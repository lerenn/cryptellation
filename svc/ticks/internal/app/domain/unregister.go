package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
)

func (t Ticks) Unregister(ctx context.Context, exchange, pair string) (int64, error) {
	count, err := t.db.DecrementSymbolListenerSubscribers(ctx, exchange, pair)
	if err != nil {
		return count, err
	}

	log := fmt.Sprintf("Unregister listener for %q on %q (count=%d)\n", exchange, pair, count)
	telemetry.L(ctx).Info(log)

	return count, nil
}
