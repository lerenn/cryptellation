package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSQLSuite(t *testing.T) {
	suite.Run(t, new(SQLSuite))
}

type SQLSuite struct {
	suite.Suite
}

func (suite *SQLSuite) TestLoadValidate() {
	cases := []struct {
		SQL SQL
		Err error
	}{
		{
			SQL: SQL{
				Host:     "host",
				Port:     1000,
				User:     "user",
				Password: "password",
				Database: "database",
			},
			Err: nil,
		},
	}

	for i, c := range cases {
		err := c.SQL.Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}
