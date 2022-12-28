package db

import (
	"errors"
	"log"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/db/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/db/sql"
)

func Init() (db.Adapter, error) {
	sqlConfig := sql.Config{}
	redisConfig := redis.Config{}

	switch {
	case sqlConfig.Load().Validate() == nil:
		log.Println("SQL Database selected")
		return sql.New()
	case redisConfig.Load().Validate() == nil:
		log.Println("Redis Database selected")
		return redis.New()
	default:
		return nil, errors.New("no DB specified")
	}
}
