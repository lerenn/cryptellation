package candlestick

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCandlestickSuite(t *testing.T) {
	suite.Run(t, new(CandlestickSuite))
}

type CandlestickSuite struct {
	suite.Suite
}

func (suite *CandlestickSuite) TestCandlestickEqual() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 1000, false}

	suite.Require().True(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualOpen() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{1, 1, 2, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualHigh() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 2, 2, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualLow() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 3, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualClose() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 4, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualVolume() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 2000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualUncomplete() {
	a := Candlestick{0, 1, 2, 3, 1000, false}
	b := Candlestick{0, 1, 2, 3, 1000, true}

	suite.Require().False(a.Equal(b))
}
