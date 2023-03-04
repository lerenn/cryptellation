package sql

import (
	"testing"

	"github.com/digital-feather/cryptellation/internal/exchanges/app/ports/db/tests"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	tests.ExchangesSuite
}

func (suite *ExchangesSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "exchanges")()

	db, err := New(config.LoadSQLConfigFromEnv())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
