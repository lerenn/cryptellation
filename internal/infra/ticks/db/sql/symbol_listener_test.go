package sql

import (
	"testing"

	"github.com/lerenn/cryptellation/internal/core/ticks/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestSymbolListenerSuite(t *testing.T) {
	suite.Run(t, new(SymbolListenerSuite))
}

type SymbolListenerSuite struct {
	db.SymbolListenerSuite
}

func (suite *SymbolListenerSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "ticks")()

	db, err := New(config.LoadSQLConfigFromEnv())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
