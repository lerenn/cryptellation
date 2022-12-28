package sql

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	tests.ExchangesSuite
}

func (suite *ExchangesSuite) SetupTest() {
	db, err := New()
	suite.Require().NoError(err)

	suite.Require().NoError(db.Reset())

	suite.DB = db
}
