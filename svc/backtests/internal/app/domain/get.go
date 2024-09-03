package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
)

func (b Backtests) Get(ctx context.Context, backtestId uuid.UUID) (backtest.Backtest, error) {
	return b.db.ReadBacktest(ctx, backtestId)
}
