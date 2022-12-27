package ticks

import (
	"context"
	"fmt"
)

func (t Ticks) Register(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	count, err := t.vdb.IncrementSymbolListenerCount(ctx, exchange, pairSymbol)
	if err != nil {
		return count, err
	}

	if count == 1 {
		err := t.launchListener(exchange, pairSymbol)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (t Ticks) launchListener(exchange, pairSymbol string) error {
	exch, exists := t.exchanges[exchange]
	if !exists {
		return fmt.Errorf("exchange %q doesn't exists", exchange)
	}

	el := internalListener{
		DB:         t.vdb,
		PubSub:     t.pubsub,
		Exchange:   exch,
		PairSymbol: pairSymbol,
	}

	return el.Run()
}
