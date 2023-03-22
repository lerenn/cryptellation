// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=backtests.go -destination=mock/backtests.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/event"
)

type Backtests interface {
	Create(ctx context.Context, payload BacktestCreationPayload) (int, error)
	ListenEvents(ctx context.Context, backtestID uint) (<-chan event.Event, error)
	Close()
}

type BacktestCreationPayload struct {
	Accounts  map[string]account.Account
	StartTime time.Time
	EndTime   *time.Time
}
