// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=client.interface.go -destination=client.mock.gen.go -package nats

package nats

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/exchange"
)

type Client interface {
	ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error)
	Close()
}
