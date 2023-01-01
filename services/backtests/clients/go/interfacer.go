// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=interfacer.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
)

type Interfacer interface {
	BacktestAccounts(ctx context.Context, backtestID uint64) (map[string]account.Account, error)
	AdvanceBacktest(ctx context.Context, backtestID uint64) error
	BacktestOrders(ctx context.Context, backtestID uint64) ([]order.Order, error)
	CreateBacktestOrder(ctx context.Context, backtestID uint64, o order.Order) error
	CreateBacktest(ctx context.Context, start, end time.Time, accounts map[string]account.Account) (id uint64, err error)
	SubscribeToBacktestEvents(ctx context.Context, backtestID uint64, exchangeName, pairSymbol string) error
	ListenBacktest(backtestID uint) (<-chan event.Event, error)
}
