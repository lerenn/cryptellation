// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/candlestick"
)

type Port interface {
	CreateCandlesticks(ctx context.Context, cs *candlestick.List) error
	ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error
	UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error
	DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error
}
