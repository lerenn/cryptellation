package exchanges

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/exchanges/deployments"
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
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "exchanges",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
