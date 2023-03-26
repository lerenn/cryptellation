// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=backtests.go -destination=mock/backtests.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/pkg/account"
	"github.com/digital-feather/cryptellation/pkg/event"
	"github.com/digital-feather/cryptellation/pkg/order"
)

type Backtests interface {
	Advance(ctx context.Context, backtestID uint) error
	Create(ctx context.Context, payload BacktestCreationPayload) (uint, error)
	CreateOrder(ctx context.Context, payload OrderCreationPayload) error
	GetAccounts(ctx context.Context, backtestID uint) (map[string]account.Account, error)
	Subscribe(ctx context.Context, backtestID uint, exchange, pair string) error
	ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error)
	Close()
}

type BacktestCreationPayload struct {
	Accounts  map[string]account.Account
	StartTime time.Time
	EndTime   *time.Time
}

type OrderCreationPayload struct {
	BacktestID   uint
	Type         order.Type
	ExchangeName string
	PairSymbol   string
	Side         order.Side
	Quantity     float64
}
