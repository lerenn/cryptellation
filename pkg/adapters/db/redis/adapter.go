package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/lerenn/cryptellation/pkg/config"
)

type Adapter struct {
	Client     *redis.Client
	ClientLock *redsync.Redsync
}

func New() (Adapter, error) {
	c := config.LoadRedis()
	if err := c.Validate(); err != nil {
		return Adapter{}, fmt.Errorf("loading redis config: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	// TODO Check connection

	pool := goredis.NewPool(client)
	ls := redsync.New(pool)

	return Adapter{
		Client:     client,
		ClientLock: ls,
	}, nil
}
