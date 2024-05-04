package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/ticks/deployments"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestSymbolListenerSuite(t *testing.T) {
	suite.Run(t, new(SymbolListenerSuite))
}

type SymbolListenerSuite struct {
	db.SymbolListenerSuite
}

func (suite *SymbolListenerSuite) SetupTest() {
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "ticks",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
