package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	db.BacktestSuite
}

func (suite *BacktestSuite) SetupTest() {
	db, err := New(context.Background(), config.LoadPostGres(nil))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.Background()))

	suite.DB = db
}
