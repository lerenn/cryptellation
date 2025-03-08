package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/db"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	db.ExchangesSuite
}

func (suite *ExchangesSuite) SetupTest() {
	db, err := New(context.Background(), config.LoadPostGres(nil))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.Background()))

	suite.DB = db
}
