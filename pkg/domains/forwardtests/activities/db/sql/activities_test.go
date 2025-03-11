package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/stretchr/testify/suite"
)

func TestForwardtestsSuite(t *testing.T) {
	suite.Run(t, new(ForwardtestsSuite))
}

type ForwardtestsSuite struct {
	db.ForwardtestSuite
}

func (suite *ForwardtestsSuite) SetupTest() {
	db, err := New(context.Background(), config.LoadPostGres(nil))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.Background()))

	suite.DB = db
}
