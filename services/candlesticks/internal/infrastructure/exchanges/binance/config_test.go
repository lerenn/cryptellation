package binance

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
		Api, Secret string
		Err         error
	}{
		{
			Api:    "api-key",
			Secret: "secret-key",
		},
	}

	var config Config
	for i, c := range cases {
		defer tmpEnvVar("BINANCE_API_KEY", c.Api)()
		defer tmpEnvVar("BINANCE_SECRET_KEY", c.Secret)()

		err := config.Load().Validate()
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
