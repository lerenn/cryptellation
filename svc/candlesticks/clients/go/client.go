// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=client.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	client "cryptellation/pkg/client"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"
)

type Client interface {
	Read(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type ReadCandlesticksPayload struct {
	Exchange string
	Pair     string
	Period   period.Symbol
	Start    *time.Time
	End      *time.Time
	Limit    uint
}
