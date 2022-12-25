package backtests

import (
	"context"
	"fmt"
)

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error {
	return b.db.LockedBacktest(backtestId, func() error {
		bt, err := b.db.ReadBacktest(ctx, backtestId)
		if err != nil {
			return fmt.Errorf("cannot get backtest: %w", err)
		}

		if _, err = bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		if err := b.db.UpdateBacktest(ctx, bt); err != nil {
			return fmt.Errorf("cannot update backtest: %w", err)
		}

		return nil
	})
}
