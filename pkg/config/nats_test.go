package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestNATSSuite(t *testing.T) {
	suite.Run(t, new(NATSSuite))
}

type NATSSuite struct {
	suite.Suite
}

func (suite *NATSSuite) TestLoadValidate() {
	cases := []struct {
		NATS NATS
		Err  error
	}{
		{
			NATS: NATS{
				Host: "host",
				Port: 1000,
			},
		},
	}

	for i, c := range cases {
		err := c.NATS.Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}
