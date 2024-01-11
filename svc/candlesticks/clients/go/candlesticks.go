// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=candlesticks.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type Client interface {
	Read(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)

	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type ReadCandlesticksPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}
