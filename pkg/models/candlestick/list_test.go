package candlestick

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

func TestCandlestickListSuite(t *testing.T) {
	suite.Run(t, new(CandlestickListSuite))
}

type CandlestickListSuite struct {
	suite.Suite
}

func (suite *CandlestickListSuite) TestNewEmpty() {
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	suite.Require().Equal("exchange", l.ExchangeName)
	suite.Require().Equal("ETH-USDC", l.PairSymbol)
	suite.Require().Equal(period.M1, l.Period)
	suite.Require().Equal(0, l.Len())
}

func (suite *CandlestickListSuite) TestNew() {
	t1 := time.Unix(0, 0)
	cs1 := Candlestick{
		Open: 1.0,
	}
	t2 := time.Unix(60, 0)
	cs2 := Candlestick{
		Open: 2.0,
	}
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(l.Set(t1, cs1))
	suite.Require().NoError(l.Set(t2, cs2))

	// Check list
	suite.Require().Equal("exchange", l.ExchangeName)
	suite.Require().Equal("ETH-USDC", l.PairSymbol)
	suite.Require().Equal(period.M1, l.Period)

	// Check candlesticks
	suite.Require().Equal(2, l.Len())

	t, e := l.Get(t1)
	suite.Require().True(e)
	suite.Require().Equal(cs1, t)

	t, e = l.Get(t2)
	suite.Require().True(e)
	suite.Require().Equal(cs2, t)
}

func (suite *CandlestickListSuite) TestNewWithUnalignedCandlestick() {
	t := time.Unix(1, 0)
	cs := Candlestick{
		Open: 1.0,
	}
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	suite.Require().Error(l.Set(t, cs))
}

func (suite *CandlestickListSuite) TestSet() {
	p := "BTC-USDC"
	csList := NewEmptyList("exchange", p, period.M1)

	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0))
	cs0bis := Candlestick{Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0bis))

	suite.Require().Equal(1, csList.Len())
	rcs0, exists := csList.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs0bis, rcs0)
}

func (suite *CandlestickListSuite) TestSetWithWrongPeriod() {
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	suite.Require().Error(l.Set(time.Unix(1, 0), cs))
}

func (suite *CandlestickListSuite) TestMerge() {
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	recvCSList := NewEmptyList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	err := recvCSList.Set(time.Unix(0, 0), cs)
	suite.Require().NoError(err)
	suite.Require().Equal(1, recvCSList.Len())

	err = l.Merge(recvCSList, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(1, l.Len())
	cs2, exists := l.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs, cs2)
}

func (suite *CandlestickListSuite) TestExtract() {
	l := NewEmptyList("exchange", "ETH-USDC", period.M1)
	for i := int64(0); i < 4; i++ {
		cs := Candlestick{
			Open:  float64(i),
			High:  0,
			Low:   0,
			Close: 0,
		}

		err := l.Set(time.Unix(60*i, 0), cs)
		suite.Require().NoError(err)
	}

	nl := l.Extract(time.Unix(60, 0), time.Unix(120, 0), 0)
	suite.Require().Equal(2, nl.Len())

	cs, exists := nl.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(1.0, cs.Open)

	cs, exists = nl.Get(time.Unix(120, 0))
	suite.Require().True(exists)
	suite.Require().Equal(2.0, cs.Open)
}

func (suite *CandlestickListSuite) TestExtractWithLimit() {
	// TODO
}

func (suite *CandlestickListSuite) TestMergeWithNotCorrespondingLists() {
	// TODO
}

func (suite *CandlestickListSuite) TestReplaceUncomplete() {
	// TODO
}

func (suite *CandlestickListSuite) TestHasUncomplete() {
	// TODO
}

func (suite *CandlestickListSuite) TestMergeIntoOne() {
	// TODO
}

func (suite *CandlestickListSuite) TestAreMissing() {
	// Given all candlesticks
	cl := NewEmptyList("exchange", "ETH-USDC", period.M1)

	for i := int64(0); i < 10; i++ {
		err := cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
		suite.Require().NoError(err)
	}

	// When asking if there is missing candlesticks
	res := cl.AreMissing(time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().False(res)
}

func (suite *CandlestickListSuite) TestAreMissingWithOneMissing() {
	// Given all candlesticks
	cl := NewEmptyList("exchange", "ETH-USDC", period.M1)

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		err := cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
		suite.Require().NoError(err)
	}

	// When asking if there is missing candlesticks
	res := cl.AreMissing(time.Unix(0, 0), time.Unix(540, 0), 0)

	// Then there is no missing
	suite.Require().True(res)
}

func (suite *CandlestickListSuite) TestAreMissingWithOneMissingAndLimit() {
	// Given all candlesticks
	cl := NewEmptyList("exchange", "ETH-USDC", period.M1)

	for i := int64(0); i < 10; i++ {
		if i == 5 {
			continue
		}

		err := cl.Set(time.Unix(60*i, 0), Candlestick{
			Open: float64(i),
		})
		suite.Require().NoError(err)
	}

	// When asking if there is missing candlesticks
	res := cl.AreMissing(time.Unix(0, 0), time.Unix(540, 0), 2)

	// Then there is no missing
	suite.Require().False(res)
}
