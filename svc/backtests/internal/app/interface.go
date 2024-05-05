package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
)

type Backtests interface {
	Advance(ctx context.Context, backtestId uuid.UUID) error
	CreateOrder(ctx context.Context, backtestId uuid.UUID, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uuid.UUID, err error)
	GetAccounts(ctx context.Context, backtestId uuid.UUID) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uuid.UUID) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uuid.UUID, exchange, pair string) error
}
