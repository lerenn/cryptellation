package sql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSqlDatabaseSuite(t *testing.T) {
	suite.Run(t, new(SqlDatabaseSuite))
}

type SqlDatabaseSuite struct {
	suite.Suite
	db *DB
}

func (suite *SqlDatabaseSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "backtests")()

	db, err := New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.db = db
}

func (suite *SqlDatabaseSuite) TestNewWithURIError() {
	defer tmpEnvVar("SQLDB_HOST", "")()

	var err error
	_, err = New()
	suite.Require().Error(err)
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
