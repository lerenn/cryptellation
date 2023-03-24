package db

import (
	"errors"

	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/services/ticks/io/db"
	"github.com/digital-feather/cryptellation/services/ticks/io/db/adapters/redis"
	"github.com/digital-feather/cryptellation/services/ticks/io/db/adapters/sql"
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
