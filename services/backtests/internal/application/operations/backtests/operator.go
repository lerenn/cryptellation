package backtests

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
)

type Operator interface {
	Advance(ctx context.Context, backtestId uint) error
	CreateOrder(ctx context.Context, backtestId uint, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uint, err error)
	GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error
}
