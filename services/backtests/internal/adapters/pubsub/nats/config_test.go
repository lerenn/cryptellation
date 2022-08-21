package nats

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
		URL string
		Err error
	}{
		{
			URL: "url",
		},
	}

	var config Config
	for i, c := range cases {
		defer tmpEnvVar("NATS_URL", c.URL)()

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
