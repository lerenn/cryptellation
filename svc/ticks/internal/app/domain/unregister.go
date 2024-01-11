package domain

import (
	"context"
	"log"
)

func (t Ticks) Unregister(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	count, err := t.db.DecrementSymbolListenerSubscribers(ctx, exchange, pairSymbol)
	if err != nil {
		return count, err
	}

	log.Printf("Unregister listener for %q on %q (count=%d)\n", exchange, pairSymbol, count)
	return count, nil
}
