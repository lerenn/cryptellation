package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
)

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uuid.UUID, exchange, pair string) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
		if _, err := bt.CreateTickSubscription(exchange, pair); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		return nil
	})
}
