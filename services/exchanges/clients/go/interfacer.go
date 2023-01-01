// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=interfacer.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

type Interfacer interface {
	ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error)
}
