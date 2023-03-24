package sql

import (
	"testing"

	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/utils"
	"github.com/digital-feather/cryptellation/services/candlesticks/io/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	tests.CandlesticksSuite
}

func (suite *CandlesticksSuite) SetupTest() {
	defer utils.TemporaryEnvVar("SQLDB_DATABASE", "candlesticks")()

	db, err := New(config.LoadSQLConfigFromEnv())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
