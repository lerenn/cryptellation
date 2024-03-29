package redis

import (
	"testing"

	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestRedisVdbSuite(t *testing.T) {
	suite.Run(t, new(RedisVdbSuite))
}

type RedisVdbSuite struct {
	db.SymbolListenerSuite
}

func (suite *RedisVdbSuite) SetupTest() {
	db, err := New()
	suite.Require().NoError(err)
	suite.DB = db
}
