package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/candlesticks/deployments"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	db.CandlesticksSuite
}

func (suite *CandlesticksSuite) SetupTest() {
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "candlesticks",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
