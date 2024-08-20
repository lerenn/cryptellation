package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/client/pkg/backtest"

	"github.com/google/uuid"
)

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uuid.UUID, exchange, pair string) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
		if _, err := bt.CreateTickSubscription(exchange, pair); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		return nil
	})
}
