// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=client.go -destination=mock.gen.go -package client

package client

import (
	"context"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type Client interface {
	Read(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}
