// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=backtests.go -destination=mock/backtests.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"
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
