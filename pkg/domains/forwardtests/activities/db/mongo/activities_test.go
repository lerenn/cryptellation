package mongo

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/stretchr/testify/suite"
)

func TestForwardTestSuite(t *testing.T) {
	suite.Run(t, new(ForwardTestSuite))
}

type ForwardTestSuite struct {
	db.ForwardTestSuite
}

func (suite *ForwardTestSuite) SetupTest() {
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
