package config

import (
	"os"
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
		Address, Password string
		Err               error
	}{
		{
			Address:  "address",
			Password: "password",
		},
	}

	for i, c := range cases {
		defer tmpEnvVar("REDIS_URL", c.Address)()
		defer tmpEnvVar("REDIS_PASSWORD", c.Password)()

		err := LoadRedisConfigFromEnv().Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
