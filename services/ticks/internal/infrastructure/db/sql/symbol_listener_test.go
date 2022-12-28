package sql

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/db/tests"
	"github.com/stretchr/testify/suite"
)

func TestSymbolListenerSuite(t *testing.T) {
	suite.Run(t, new(SymbolListenerSuite))
}

type SymbolListenerSuite struct {
	tests.SymbolListenerSuite
}

func (suite *SymbolListenerSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "ticks")()

	db, err := New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
