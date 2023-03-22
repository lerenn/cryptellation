package app

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/types/order"
)

func (b Backtests) GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}
