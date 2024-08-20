// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
)

type Port interface {
	CreateCandlesticks(ctx context.Context, cs *candlestick.List) error
	ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error
	UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error
	DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error
}
