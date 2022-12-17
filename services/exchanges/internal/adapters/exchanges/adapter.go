// Generate code for mock
//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

type Adapter interface {
	Infos(ctx context.Context) (exchange.Exchange, error)
}
