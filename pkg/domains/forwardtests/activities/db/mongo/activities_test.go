package mongo

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/stretchr/testify/suite"
)

func TestForwardtestSuite(t *testing.T) {
	suite.Run(t, new(ForwardtestSuite))
}

type ForwardtestSuite struct {
	db.ForwardtestSuite
}

func (suite *ForwardtestSuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(&config.Mongo{
			Database: "cryptellation-forwardtests-integration-tests",
		}),
	)
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.Background()))

	suite.DB = db
}
