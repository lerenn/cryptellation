package mongo

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/candlesticks/internal/app/ports/db"

	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	db.CandlesticksSuite
}

func (suite *CandlesticksSuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(
			&config.Mongo{
				Database: "cryptellation-candlesticks-integration-tests",
			}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
