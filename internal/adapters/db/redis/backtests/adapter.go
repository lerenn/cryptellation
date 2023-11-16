package backtests

import (
	adapter "github.com/lerenn/cryptellation/internal/adapters/db/redis"
)

type Adapter struct {
	redis adapter.Adapter
}

func New() (Adapter, error) {
	// Create embedded database access
	db, err := adapter.New()

	// Return database access
	return Adapter{
		redis: db,
	}, err
}
