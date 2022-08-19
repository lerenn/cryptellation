package cockroach

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) TestLoadValidate() {
	cases := []struct {
		Host, Port, User, Password, Database string
		Err                                  error
	}{
		{
			Host:     "host",
			Port:     "1000",
			User:     "user",
			Password: "password",
			Database: "database",
			Err:      nil,
		},
	}

	var config Config
	for i, c := range cases {
		defer tmpEnvVar("COCKROACHDB_HOST", c.Host)()
		defer tmpEnvVar("COCKROACHDB_PORT", c.Port)()
		defer tmpEnvVar("COCKROACHDB_USER", c.User)()
		defer tmpEnvVar("COCKROACHDB_PASSWORD", c.Password)()
		defer tmpEnvVar("COCKROACHDB_DATABASE", c.Database)()

		err := config.Load().Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}
