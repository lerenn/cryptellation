package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/configs/sql/down"
	"github.com/lerenn/cryptellation/v1/configs/sql/up"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/migrator"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	db.ExchangesSuite
}

func (suite *ExchangesSuite) SetupSuite() {
	db, err := New(context.Background(), config.LoadPostGres(nil))
	suite.Require().NoError(err)

	mig, err := migrator.NewMigrator(context.Background(), db.client.DB, up.Migrations, down.Migrations, nil)
	suite.Require().NoError(err)
	suite.Require().NoError(mig.MigrateToLatest(context.Background()))

	suite.DB = db
}

func (suite *ExchangesSuite) SetupTest() {
	db := suite.DB.(*Activities)
	suite.Require().NoError(db.Reset(context.Background()))
}
