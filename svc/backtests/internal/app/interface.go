package app

import (
	"context"

	"cryptellation/pkg/models/account"
	"cryptellation/pkg/models/order"

	"cryptellation/svc/backtests/pkg/backtest"

	"github.com/google/uuid"
)

type Backtests interface {
	Advance(ctx context.Context, backtestId uuid.UUID) error
	CreateOrder(ctx context.Context, backtestId uuid.UUID, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uuid.UUID, err error)
	GetAccounts(ctx context.Context, backtestId uuid.UUID) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uuid.UUID) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uuid.UUID, exchange, pair string) error
}
