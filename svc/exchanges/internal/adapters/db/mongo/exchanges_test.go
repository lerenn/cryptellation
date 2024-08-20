package mongo

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/db"

	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	db.ExchangesSuite
}

func (suite *ExchangesSuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(
			&config.Mongo{
				Database: "cryptellation-exchanges-integration-tests",
			}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
