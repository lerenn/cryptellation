package backtests

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/order"
)

func (b Backtests) GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Orders, nil
}
