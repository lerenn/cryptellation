package backtests

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
)

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) error {
		if _, err := bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		return nil
	})
}
