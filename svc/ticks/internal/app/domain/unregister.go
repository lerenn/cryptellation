package domain

import (
	"context"
	"log"
)

func (t Ticks) Unregister(ctx context.Context, exchange, pair string) (int64, error) {
	count, err := t.db.DecrementSymbolListenerSubscribers(ctx, exchange, pair)
	if err != nil {
		return count, err
	}

	log.Printf("Unregister listener for %q on %q (count=%d)\n", exchange, pair, count)
	return count, nil
}
