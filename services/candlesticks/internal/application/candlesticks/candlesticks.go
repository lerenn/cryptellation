package candlesticks

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
)

// Test interface implementation
var _ Operator = (*Candlesticks)(nil)

type Candlesticks struct {
	repository db.Adapter
	services   map[string]exchanges.Adapter
}

func New(repository db.Adapter, services map[string]exchanges.Adapter) Candlesticks {
	if repository == nil {
		panic("nil repository")
	}

	if len(services) == 0 {
		panic("nil services")
	}

	return Candlesticks{
		repository: repository,
		services:   services,
	}
}
