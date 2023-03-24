// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/pkg/backtest"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type LockedBacktestCallback func(bt *backtest.Backtest) error

type Port interface {
	CreateBacktest(ctx context.Context, bt *backtest.Backtest) error
	ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error)
	UpdateBacktest(ctx context.Context, bt backtest.Backtest) error
	DeleteBacktest(ctx context.Context, bt backtest.Backtest) error

	LockedBacktest(ctx context.Context, id uint, fn LockedBacktestCallback) error
}
