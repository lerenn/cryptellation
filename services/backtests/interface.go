package backtests

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/account"
	"github.com/lerenn/cryptellation/pkg/backtest"
	"github.com/lerenn/cryptellation/pkg/order"
)

type Interface interface {
	Advance(ctx context.Context, backtestId uint) error
	CreateOrder(ctx context.Context, backtestId uint, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uint, err error)
	GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error
}
