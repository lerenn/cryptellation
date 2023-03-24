// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=candlesticks.go -destination=mock/candlesticks.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
)

type Candlesticks interface {
	Read(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)
	Close()
}

type ReadCandlesticksPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}
