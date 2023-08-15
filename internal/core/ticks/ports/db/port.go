// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"
)

const (
	Expiration = 3 * time.Second
	RetryDelay = 100 * time.Millisecond
	Tries      = 20
)

type Port interface {
	// Symbol Listener Count
	IncrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error)
	DecrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error)
	GetSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error)
	ClearSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) error
	ClearAllSymbolListenersCount(ctx context.Context) error
}
