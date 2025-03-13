package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/configs/sql/down"
	"github.com/lerenn/cryptellation/v1/configs/sql/up"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/migrator"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	db.BacktestSuite
}

func (suite *BacktestSuite) SetupSuite() {
	db, err := New(context.Background(), config.LoadSQL(nil))
	suite.Require().NoError(err)

	mig, err := migrator.NewMigrator(context.Background(), db.client.DB, up.Migrations, down.Migrations, nil)
	suite.Require().NoError(err)
	suite.Require().NoError(mig.MigrateToLatest(context.Background()))

	suite.DB = db
}

func (suite *BacktestSuite) SetupTest() {
	db := suite.DB.(*Activities)
	suite.Require().NoError(db.Reset(context.Background()))
}
