// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=backtests.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
)

type Client interface {
	Advance(ctx context.Context, backtestID uuid.UUID) error
	Create(ctx context.Context, payload BacktestCreationPayload) (uuid.UUID, error)
	CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error
	GetAccounts(ctx context.Context, backtestID uuid.UUID) (map[string]account.Account, error)
	Subscribe(ctx context.Context, backtestID uuid.UUID, exchange, pair string) error
	ListenEvents(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type BacktestCreationPayload struct {
	Accounts  map[string]account.Account
	StartTime time.Time
	EndTime   *time.Time
}
