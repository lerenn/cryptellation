package candlesticks

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/internal/components/candlesticks/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	db.CandlesticksSuite
}

func (suite *CandlesticksSuite) SetupTest() {
	defer utils.TemporaryEnvVar("SQLDB_DATABASE", "candlesticks")()

	db, err := New(config.LoadSQL())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
