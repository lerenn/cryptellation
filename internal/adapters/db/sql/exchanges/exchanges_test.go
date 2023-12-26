package exchanges

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/internal/components/exchanges/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	db.ExchangesSuite
}

func (suite *ExchangesSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "exchanges")()

	db, err := New(config.LoadSQL())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
