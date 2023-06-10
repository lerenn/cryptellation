// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=exchanges.go -destination=mock/exchanges.gen.go -package mock

package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/exchange"
)

type Exchanges interface {
	Read(ctx context.Context, names ...string) ([]exchange.Exchange, error)
	Close()
}
