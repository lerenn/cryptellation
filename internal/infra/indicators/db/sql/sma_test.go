package sql

import (
	"testing"

	"github.com/lerenn/cryptellation/internal/core/indicators/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestIndicatorsSuite(t *testing.T) {
	suite.Run(t, new(IndicatorsSuite))
}

type IndicatorsSuite struct {
	db.IndicatorsSuite
}

func (suite *IndicatorsSuite) SetupTest() {
	defer utils.TemporaryEnvVar("SQLDB_DATABASE", "indicators")()

	db, err := New(config.LoadSQLConfigFromEnv())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.DB = db
}
