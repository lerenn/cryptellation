package exchanges

import (
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/exchanges"
)

type Exchanges struct {
	db       db.Adapter
	services map[string]exchanges.Adapter
}

func New(db db.Adapter, services map[string]exchanges.Adapter) Exchanges {
	if db == nil {
		panic("nil db")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return Exchanges{
		db:       db,
		services: services,
	}
}
