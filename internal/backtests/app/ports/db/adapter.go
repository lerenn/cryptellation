// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type LockedBacktestCallback func(bt *domain.Backtest) error

type Adapter interface {
	CreateBacktest(ctx context.Context, bt *domain.Backtest) error
	ReadBacktest(ctx context.Context, id uint) (domain.Backtest, error)
	UpdateBacktest(ctx context.Context, bt domain.Backtest) error
	DeleteBacktest(ctx context.Context, bt domain.Backtest) error

	LockedBacktest(ctx context.Context, id uint, fn LockedBacktestCallback) error
}
