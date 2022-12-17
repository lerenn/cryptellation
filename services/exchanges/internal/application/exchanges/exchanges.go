package exchanges

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/adapters/exchanges"
)

type Exchanges struct {
	repository db.Adapter
	services   map[string]exchanges.Adapter
}

func New(repository db.Adapter, services map[string]exchanges.Adapter) Exchanges {
	if repository == nil {
		panic("nil repository")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return Exchanges{
		repository: repository,
		services:   services,
	}
}
