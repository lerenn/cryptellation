package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBinanceSuite(t *testing.T) {
	suite.Run(t, new(BinanceSuite))
}

type BinanceSuite struct {
	suite.Suite
}

func (suite *BinanceSuite) TestLoadValidate() {
	cases := []struct {
		Binance Binance
		Err     error
	}{
		{
			Binance: Binance{
				ApiKey:    "api-key",
				SecretKey: "secret-key",
			},
		},
	}

	for i, c := range cases {
		err := c.Binance.Validate()
		suite.Require().Equal(c.Err, err, i)
	}
}
