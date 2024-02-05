package domain

import (
	"context"
	"log"
)

func (t Ticks) Register(ctx context.Context, exchange, pair string) (int64, error) {
	count, err := t.db.IncrementSymbolListenerSubscribers(ctx, exchange, pair)
	if err != nil {
		return count, err
	}

	if count == 1 {
		err := t.launchListener(exchange, pair)
		if err != nil {
			return count, err
		}
	}

	log.Printf("Register listener for %q on %q (count=%d)\n", exchange, pair, count)
	return count, nil
}

func (t Ticks) launchListener(exchange, pair string) error {

	el := internalListener{
		DB:        t.db,
		Events:    t.events,
		Exchanges: t.exchanges,

		Exchange: exchange,
		Pair:     pair,
	}

	return el.Run()
}
