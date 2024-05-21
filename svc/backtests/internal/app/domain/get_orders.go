package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/order"
)

func (b Backtests) GetOrders(ctx context.Context, backtestId uuid.UUID) ([]order.Order, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}
