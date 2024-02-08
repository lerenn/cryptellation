package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
)

func (t Ticks) Register(ctx context.Context, exchange, pair string) (int64, error) {
	count, err := t.db.IncrementSymbolListenerSubscribers(ctx, exchange, pair)
	if err != nil {
		return count, err
	}

	if count == 1 {
		err := t.launchListener(ctx, exchange, pair)
		if err != nil {
			return count, err
		}
	}

	telemetry.L(ctx).Info(fmt.Sprintf("Register listener for %q on %q (count=%d)\n", exchange, pair, count))
	return count, nil
}

func (t Ticks) launchListener(ctx context.Context, exchange, pair string) error {

	el := internalListener{
		DB:        t.db,
		Events:    t.events,
		Exchanges: t.exchanges,

		Exchange: exchange,
		Pair:     pair,
	}

	return el.Run(ctx)
}
