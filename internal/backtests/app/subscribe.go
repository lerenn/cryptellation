package app

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
)

func (b Backtests) SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *domain.Backtest) error {
		if _, err := bt.CreateTickSubscription(exchange, pairSymbol); err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		return nil
	})
}
