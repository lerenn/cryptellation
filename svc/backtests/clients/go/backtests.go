// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=backtests.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
)

type Client interface {
	Advance(ctx context.Context, backtestID uuid.UUID) error
	Create(ctx context.Context, payload BacktestCreationPayload) (uuid.UUID, error)
	CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error
	Get(ctx context.Context, backtestID uuid.UUID) (backtest.Backtest, error)
	GetAccounts(ctx context.Context, backtestID uuid.UUID) (map[string]account.Account, error)
	Subscribe(ctx context.Context, backtestID uuid.UUID, exchange, pair string) error
	ListenEvents(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error)
	List(ctx context.Context) ([]backtest.Backtest, error)
	ListOrders(ctx context.Context, backtestID uuid.UUID) ([]order.Order, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type BacktestCreationPayload struct {
	Accounts    map[string]account.Account
	StartTime   time.Time
	EndTime     *time.Time
	Mode        *backtest.Mode
	PricePeriod *period.Symbol
}
