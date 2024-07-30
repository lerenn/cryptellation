package domain

import (
	"context"

	"cryptellation/pkg/models/order"

	"github.com/google/uuid"
)

func (b Backtests) GetOrders(ctx context.Context, backtestId uuid.UUID) ([]order.Order, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}
