package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
)

func (b Backtests) List(ctx context.Context) ([]backtest.Backtest, error) {
	bts, err := b.db.ListBacktests(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when listing backtests: %w", err)
	}

	return bts, nil
}
