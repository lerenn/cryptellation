// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
)

type Port interface {
	CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error
	ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error)
	UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error
	DeleteExchanges(ctx context.Context, names ...string) error
}
