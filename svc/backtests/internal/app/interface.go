package app

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/order"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/google/uuid"
)

type Backtests interface {
	Advance(ctx context.Context, backtestId uuid.UUID) error
	CreateOrder(ctx context.Context, backtestId uuid.UUID, order order.Order) error
	Create(ctx context.Context, req backtest.NewPayload) (id uuid.UUID, err error)
	Get(ctx context.Context, backtestId uuid.UUID) (backtest.Backtest, error)
	GetAccounts(ctx context.Context, backtestId uuid.UUID) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uuid.UUID) ([]order.Order, error)
	List(ctx context.Context) ([]backtest.Backtest, error)
	SubscribeToEvents(ctx context.Context, backtestId uuid.UUID, exchange, pair string) error
}
