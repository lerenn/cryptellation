package redis

import (
	"fmt"

	adapter "github.com/lerenn/cryptellation/pkg/adapters/db/redis"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/ticks/deployments"
)

type Adapter struct {
	redis adapter.Adapter
}

func New() (*Adapter, error) {
	// Create embedded database access
	db, err := adapter.New(config.LoadRedis(
		&config.Redis{
			Address: fmt.Sprintf("localhost:%d", deployments.DockerComposeRedisPort),
		},
	))

	// Return database access
	return &Adapter{
		redis: db,
	}, err
}
