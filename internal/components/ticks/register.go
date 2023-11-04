package ticks

import (
	"context"
	"log"
)

func (t Ticks) Register(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	count, err := t.db.IncrementSymbolListenerSubscribers(ctx, exchange, pairSymbol)
	if err != nil {
		return count, err
	}

	if count == 1 {
		err := t.launchListener(exchange, pairSymbol)
		if err != nil {
			return count, err
		}
	}

	log.Printf("Register listener for %q on %q (count=%d)\n", exchange, pairSymbol, count)
	return count, nil
}

func (t Ticks) launchListener(exchange, pairSymbol string) error {

	el := internalListener{
		DB:        t.db,
		Events:    t.events,
		Exchanges: t.exchanges,

		ExchangeName: exchange,
		PairSymbol:   pairSymbol,
	}

	return el.Run()
}
