package redis

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestRedisVdbSuite(t *testing.T) {
	suite.Run(t, new(RedisVdbSuite))
}

type RedisVdbSuite struct {
	tests.SymbolListenerSuite
}

func (suite *RedisVdbSuite) SetupTest() {
	db, err := New()
	suite.Require().NoError(err)
	suite.DB = db
}
