package app

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
)

func (b Backtests) Create(ctx context.Context, req domain.NewPayload) (id uint, err error) {
	bt, err := domain.New(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("creating a new backtest from request: %w", err)
	}

	err = b.db.CreateBacktest(ctx, &bt)
	if err != nil {
		return 0, fmt.Errorf("adding backtest to vdb: %w", err)
	}

	return bt.ID, nil
}
