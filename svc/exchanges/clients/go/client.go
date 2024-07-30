// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=client.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"cryptellation/pkg/client"

	"cryptellation/svc/exchanges/pkg/exchange"
)

type Client interface {
	Read(ctx context.Context, names ...string) ([]exchange.Exchange, error)
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}
