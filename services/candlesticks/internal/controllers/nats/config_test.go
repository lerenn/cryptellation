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
		Host, Port string
		Err        error
	}{
		{
			Host: "host",
			Port: "1000",
		},
	}

	for i, c := range cases {
		defer tmpEnvVar("NATS_HOST", c.Host)()
		defer tmpEnvVar("NATS_PORT", c.Port)()

		err := loadConfig().Validate()
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
