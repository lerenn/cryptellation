package db

import (
	"errors"

	"github.com/lerenn/cryptellation/internal/core/ticks/io/db"
	"github.com/lerenn/cryptellation/internal/core/ticks/io/db/adapters/redis"
	"github.com/lerenn/cryptellation/internal/core/ticks/io/db/adapters/sql"
	"github.com/lerenn/cryptellation/pkg/config"
)

var (
	ErrNoValidDatabase = errors.New("no valid database")
)

func New(sqlCfg config.SQL, redisCfg config.Redis) (db.Port, error) {
	switch {
	case sqlCfg.Validate() == nil:
		return sql.New(sqlCfg)
	case redisCfg.Validate() == nil:
		return redis.New(redisCfg)
	default:
		return nil, ErrNoValidDatabase
	}
}
