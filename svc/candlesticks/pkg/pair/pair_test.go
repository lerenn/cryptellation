package pair

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPairsSuite(t *testing.T) {
	suite.Run(t, new(PairsSuite))
}

type PairsSuite struct {
	suite.Suite
}

func (suite *PairsSuite) TestFormatPair() {
	cases := []struct {
		BaseSymbol  string
		QuoteSymbol string
		Pair        string
	}{
		{
			BaseSymbol:  "ETH",
			QuoteSymbol: "BTC",
			Pair:        "ETH-BTC",
		},
	}

	for i, c := range cases {
		symbol := FormatPair(c.BaseSymbol, c.QuoteSymbol)
		suite.Require().Equal(c.Pair, symbol, i)
	}
}

func (suite *PairsSuite) TestParsePair() {
	cases := []struct {
		Pair        string
		BaseSymbol  string
		QuoteSymbol string
		Error       bool
	}{
		{
			Pair:        "ETH-BTC",
			BaseSymbol:  "ETH",
			QuoteSymbol: "BTC",
			Error:       false,
		}, {
			Pair:        "",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       true,
		}, {
			Pair:        "-",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       false,
		}, {
			Pair:        "--",
			BaseSymbol:  "",
			QuoteSymbol: "",
			Error:       true,
		},
	}

	for i, c := range cases {
		baseSymbol, quoteSymbol, err := ParsePair(c.Pair)
		suite.Require().Equal(c.BaseSymbol, baseSymbol, i)
		suite.Require().Equal(c.QuoteSymbol, quoteSymbol, i)
		if c.Error {
			suite.Require().Error(err, i)
		} else {
			suite.Require().NoError(err, i)
		}
	}
}
