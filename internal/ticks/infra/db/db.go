package db

import (
	"errors"

	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/ticks/infra/db/redis"
	"github.com/digital-feather/cryptellation/internal/ticks/infra/db/sql"
	"github.com/digital-feather/cryptellation/pkg/config"
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
