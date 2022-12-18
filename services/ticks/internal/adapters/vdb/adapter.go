// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package vdb

package vdb

import (
	"context"
	"time"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type Adapter interface {
	IncrementSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	DecrementSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	GetSymbolListenerCount(ctx context.Context, exchange, pairSymbol string) (int64, error)
	ClearSymbolListenersCount(ctx context.Context) error
}
