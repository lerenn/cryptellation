package backtests

import (
	adapter "github.com/lerenn/cryptellation/pkg/adapters/db/redis"
	"github.com/lerenn/cryptellation/pkg/config"
)

type Adapter struct {
	redis adapter.Adapter
}

func New(c config.Redis) (Adapter, error) {
	// Create embedded database access
	db, err := adapter.New(c)

	// Return database access
	return Adapter{
		redis: db,
	}, err
}
