package binance

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
		Config Config
		Err    error
	}{
		{
			Config: Config{
				ApiKey:    "api-key",
				SecretKey: "secret-key",
			},
		},
	}

	for i, c := range cases {
		err := c.Config.Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}
