package daemon

import (
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/internal/app"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/domain"
)

type components struct {
	exchanges exchanges.Exchanges
}

func newComponents(adapters adapters) components {
	return components{
		exchanges: domain.New(adapters.db, adapters.exchanges),
	}
}

func (c components) Close() {
}
