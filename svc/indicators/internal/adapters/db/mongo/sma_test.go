package mongo

import (
	"context"
	"testing"

	"cryptellation/internal/config"

	"cryptellation/svc/indicators/internal/app/ports/db"

	"github.com/stretchr/testify/suite"
)

func TestSMASuite(t *testing.T) {
	suite.Run(t, new(SMASuite))
}

type SMASuite struct {
	db.IndicatorsSuite
}

func (suite *SMASuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(
			&config.Mongo{
				Database: "cryptellation-indicators-integration-tests",
			}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
