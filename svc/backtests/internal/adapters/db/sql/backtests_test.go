package backtests

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/backtests/deployments"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	db.BacktestSuite
}

func (suite *BacktestSuite) SetupTest() {
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "backtests",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
