// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=backtests.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
)

type Client interface {
	Advance(ctx context.Context, backtestID uint) error
	Create(ctx context.Context, payload BacktestCreationPayload) (uint, error)
	CreateOrder(ctx context.Context, payload OrderCreationPayload) error
	GetAccounts(ctx context.Context, backtestID uint) (map[string]account.Account, error)
	Subscribe(ctx context.Context, backtestID uint, exchange, pair string) error
	ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
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
