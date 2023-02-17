package sql

import (
	"os"
	"strconv"
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
	db, err := newTestDB()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.db = db
}

func (suite *SqlDatabaseSuite) TestNewWithURIError() {
	var err error
	_, err = New(Config{})
	suite.Require().Error(err)
}

func newTestDB() (*DB, error) {

	port, _ := strconv.Atoi(os.Getenv("SQLDB_PORT"))
	return New(Config{
		Host:     os.Getenv("SQLDB_HOST"),
		Port:     port,
		User:     os.Getenv("SQLDB_USER"),
		Password: os.Getenv("SQLDB_PASSWORD"),
		Database: "candlesticks",
	})
}
