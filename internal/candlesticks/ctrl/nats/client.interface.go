// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=client.interface.go -destination=client.mock.gen.go -package nats

package nats

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
)

type Client interface {
	ReadCandlesticks(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)
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
