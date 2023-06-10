package redis

import (
	"testing"

	"github.com/lerenn/cryptellation/internal/core/ticks/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestRedisVdbSuite(t *testing.T) {
	suite.Run(t, new(RedisVdbSuite))
}

type RedisVdbSuite struct {
	db.SymbolListenerSuite
}

func (suite *RedisVdbSuite) SetupTest() {
	db, err := New(config.LoadRedisConfigFromEnv())
	suite.Require().NoError(err)
	suite.DB = db
}
