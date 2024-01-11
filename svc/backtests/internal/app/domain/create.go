package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
)

func (b Backtests) Create(ctx context.Context, req backtest.NewPayload) (id uint, err error) {
	bt, err := backtest.New(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("creating a new backtest from request: %w", err)
	}

	err = b.db.CreateBacktest(ctx, &bt)
	if err != nil {
		return 0, fmt.Errorf("adding backtest to db: %w", err)
	}

	return bt.ID, nil
}
