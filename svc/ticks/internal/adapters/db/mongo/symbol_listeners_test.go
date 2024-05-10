package mongo

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestSymbolListenersSuite(t *testing.T) {
	suite.Run(t, new(SymbolListenersSuite))
}

type SymbolListenersSuite struct {
	db.SymbolListenerSuite
}

func (suite *SymbolListenersSuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(
			&config.Mongo{
				Database: "cryptellation-ticks-integration-tests",
			}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
