// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=mock/exchanges.gen.go -package mock

package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/exchange"
)

type Exchanges interface {
	Read(ctx context.Context, names ...string) ([]exchange.Exchange, error)

	ServiceInfo(ctx context.Context) (ServiceInfo, error)
	Close(ctx context.Context)
}
