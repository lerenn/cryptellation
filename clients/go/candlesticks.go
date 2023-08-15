// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=candlesticks.go -destination=mock/candlesticks.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
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
