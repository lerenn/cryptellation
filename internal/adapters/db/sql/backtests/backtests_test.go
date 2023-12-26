package backtests

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/internal/components/backtests/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	db.BacktestSuite
}

func (suite *BacktestSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "backtests")()

	db, err := New(config.LoadSQL())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
