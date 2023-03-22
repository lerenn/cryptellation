package app

import (
	"context"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/order"
)

type Controller interface {
	Advance(ctx context.Context, backtestId uint) error
	CreateOrder(ctx context.Context, backtestId uint, order order.Order) error
	Create(ctx context.Context, req domain.NewPayload) (id uint, err error)
	GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error)
	GetOrders(ctx context.Context, backtestId uint) ([]order.Order, error)
	SubscribeToEvents(ctx context.Context, backtestId uint, exchange, pairSymbol string) error
}
