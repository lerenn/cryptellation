package daemon

import (
	exchanges "cryptellation/svc/exchanges/internal/app"
	"cryptellation/svc/exchanges/internal/app/domain"
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
