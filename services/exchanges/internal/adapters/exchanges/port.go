package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

type Port interface {
	Infos(ctx context.Context) (exchange.Exchange, error)
}
