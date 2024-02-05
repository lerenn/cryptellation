package app

import (
	"context"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
)

type Backtests interface {
	Advance(ctx context.Context, backtestId uint) error
	CreateOrder(ctx context.Context, backtestId uint, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uint, err error)
	GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pair string) error
}
