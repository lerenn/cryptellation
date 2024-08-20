package domain

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/google/uuid"
)

func (b Backtests) Create(ctx context.Context, req backtest.NewPayload) (id uuid.UUID, err error) {
	bt, err := backtest.New(ctx, req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("creating a new backtest from request: %w", err)
	}

	err = b.db.CreateBacktest(ctx, bt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("adding backtest to db: %w", err)
	}

	return bt.ID, nil
}
