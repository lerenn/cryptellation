package candlestick

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestCandlestickSuite(t *testing.T) {
	suite.Run(t, new(CandlestickSuite))
}

type CandlestickSuite struct {
	suite.Suite
}

func (suite *CandlestickSuite) TestCandlestickEqual() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}

	suite.Require().True(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualOpen() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 1, 1, 2, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualHigh() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 2, 2, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualLow() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 1, 3, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualClose() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 1, 2, 4, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualVolume() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 2000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualUncomplete() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, true}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickNotEqualTime() {
	a := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}
	b := Candlestick{time.Unix(60, 0), 0, 1, 2, 3, 1000, false}

	suite.Require().False(a.Equal(b))
}

func (suite *CandlestickSuite) TestCandlestickPrice() {
	c := Candlestick{time.Unix(0, 0), 0, 1, 2, 3, 1000, false}

	v := c.Price(PriceIsOpen)
	if v != 0 {
		suite.Require().FailNow("Wrong value:", v)
	}

	v = c.Price(PriceIsHigh)
	if v != 1 {
		suite.Require().FailNow("Wrong value:", v)
	}

	v = c.Price(PriceIsLow)
	if v != 2 {
		suite.Require().FailNow("Wrong value:", v)
	}

	v = c.Price(PriceIsClose)
	if v != 3 {
		suite.Require().FailNow("Wrong value:", v)
	}
}
