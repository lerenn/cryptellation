package candlesticks

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/exchanges"
)

// Test interface implementation
var _ Operator = (*Candlesticks)(nil)

type Candlesticks struct {
	db       db.Adapter
	services map[string]exchanges.Adapter
}

func New(db db.Adapter, services map[string]exchanges.Adapter) Candlesticks {
	if db == nil {
		panic("nil db")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return Candlesticks{
		db:       db,
		services: services,
	}
}
