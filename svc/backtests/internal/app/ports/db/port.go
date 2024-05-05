// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type LockedBacktestCallback func(bt *backtest.Backtest) error

type Port interface {
	CreateBacktest(ctx context.Context, bt backtest.Backtest) error
	ReadBacktest(ctx context.Context, id uuid.UUID) (backtest.Backtest, error)
	UpdateBacktest(ctx context.Context, bt backtest.Backtest) error
	DeleteBacktest(ctx context.Context, bt backtest.Backtest) error

	LockedBacktest(ctx context.Context, id uuid.UUID, fn LockedBacktestCallback) error
}
