package sql

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	tests.CandlesticksSuite
}

func (suite *CandlesticksSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "candlesticks")()

	db, err := New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
