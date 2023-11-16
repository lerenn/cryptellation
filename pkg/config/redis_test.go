package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

type RedisSuite struct {
	suite.Suite
}

func (suite *RedisSuite) TestLoadValidate() {
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

		cfg := LoadRedis()
		suite.Require().Equal(c.Err, cfg.Validate(), i)
	}
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
