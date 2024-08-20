package candlestick

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/stretchr/testify/suite"
)

func TestCandlestickListSuite(t *testing.T) {
	suite.Run(t, new(CandlestickListSuite))
}

type CandlestickListSuite struct {
	suite.Suite
}

func (suite *CandlestickListSuite) TestNewEmpty() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().Equal("exchange", l.Exchange)
	suite.Require().Equal("ETH-USDC", l.Pair)
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
	l := NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(l.Set(t1, cs1))
	suite.Require().NoError(l.Set(t2, cs2))

	// Check list
	suite.Require().Equal("exchange", l.Exchange)
	suite.Require().Equal("ETH-USDC", l.Pair)
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
	l := NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().Error(l.Set(t, cs))
}

func (suite *CandlestickListSuite) TestMustSet() {
	// TODO: set
}

func (suite *CandlestickListSuite) TestSet() {
	p := "BTC-USDC"
	csList := NewList("exchange", p, period.M1)

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
	l := NewList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	suite.Require().Error(l.Set(time.Unix(1, 0), cs))
}

func (suite *CandlestickListSuite) TestMerge() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	recvCSList := NewList("exchange", "ETH-USDC", period.M1)
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
	l := NewList("exchange", "ETH-USDC", period.M1)
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

func (suite *CandlestickListSuite) TestFillMissing() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(60, 0), cs0))

	// Fill missing candlesticks
	err := csList.FillMissing(time.Unix(0, 0), time.Unix(180, 0), Candlestick{Open: 10, High: 20, Low: 5, Close: 15})
	suite.Require().NoError(err)

	// Check that the existing one is not overwritten
	cs, exists := csList.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}, cs)

	// Check that missing candlesticks has been filled
	for _, m := range []int64{0, 2, 3} {
		cs, exists := csList.Get(time.Unix(m*60, 0))
		suite.Require().True(exists)
		suite.Require().Equal(Candlestick{Open: 10, High: 20, Low: 5, Close: 15}, cs)
	}
}

func (suite *CandlestickListSuite) TestGetUncompleteTimes() {
	// TODO
}

func (suite *CandlestickListSuite) TestGetMissingTimes() {

}
