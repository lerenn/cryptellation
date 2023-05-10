package redis

import (
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/services/ticks/io/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestRedisVdbSuite(t *testing.T) {
	suite.Run(t, new(RedisVdbSuite))
}

type RedisVdbSuite struct {
	tests.SymbolListenerSuite
}

func (suite *RedisVdbSuite) SetupTest() {
	db, err := New(config.LoadRedisConfigFromEnv())
	suite.Require().NoError(err)
	suite.DB = db
}
