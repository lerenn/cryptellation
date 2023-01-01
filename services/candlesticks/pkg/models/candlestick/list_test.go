package candlestick

import (
	"errors"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/timeserie"
	"github.com/stretchr/testify/suite"
)

func TestCandlestickListSuite(t *testing.T) {
	suite.Run(t, new(CandlestickListSuite))
}

type CandlestickListSuite struct {
	suite.Suite
}

func (suite *CandlestickListSuite) TestNewEmpty() {
	id := ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	}
	l := NewEmptyList(id)
	suite.Require().Equal("exchange", l.ExchangeName())
	suite.Require().Equal("ETH-USDC", l.PairSymbol())
	suite.Require().Equal(period.M1, l.Period())
	suite.Require().Equal(0, l.Len())
	suite.Require().Equal(id, l.ID())
}

func (suite *CandlestickListSuite) TestNew() {
	id := ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	}
	tc1 := TimedCandlestick{Time: time.Unix(0, 0), Candlestick: Candlestick{
		Open: 1.0,
	}}
	tc2 := TimedCandlestick{Time: time.Unix(60, 0), Candlestick: Candlestick{
		Open: 2.0,
	}}
	l, err := NewList(id, tc1, tc2)

	// Check list
	suite.Require().NoError(err)
	suite.Require().Equal("exchange", l.ExchangeName())
	suite.Require().Equal("ETH-USDC", l.PairSymbol())
	suite.Require().Equal(period.M1, l.Period())
	suite.Require().Equal(id, l.ID())

	// Check candlesticks
	suite.Require().Equal(2, l.Len())

	t, e := l.Get(tc1.Time)
	suite.Require().True(e)
	suite.Require().Equal(tc1.Candlestick, t)

	t, e = l.Get(tc2.Time)
	suite.Require().True(e)
	suite.Require().Equal(tc2.Candlestick, t)
}

func (suite *CandlestickListSuite) TestNewWithUnalignedCandlestick() {
	id := ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	}
	tc := TimedCandlestick{Time: time.Unix(1, 0), Candlestick: Candlestick{
		Open: 1.0,
	}}
	_, err := NewList(id, tc)

	// Check list
	suite.Require().Error(err)
}

func (suite *CandlestickListSuite) TestMergeTimeSeries() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	ts := timeserie.New()
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	ts.Set(time.Unix(0, 0), cs)
	suite.Require().Equal(1, ts.Len())

	err := l.MergeTimeSeries(*ts, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(1, l.Len())
	cs2, exists := l.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs, cs2)
}

func (suite *CandlestickListSuite) TestSetWrongPeriod() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	suite.Require().Error(l.Set(time.Unix(1, 0), cs))
}

func (suite *CandlestickListSuite) TestMerge() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	recvCSList := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	err := recvCSList.Set(time.Unix(0, 0), cs)
	suite.Require().NoError(err)
	suite.Require().Equal(1, recvCSList.Len())

	err = l.Merge(*recvCSList, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(1, l.Len())
	cs2, exists := l.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs, cs2)
}

func (suite *CandlestickListSuite) TestDelete() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	cs := Candlestick{
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	err := l.Set(time.Unix(0, 0), cs)
	suite.Require().NoError(err)
	suite.Require().Equal(1, l.Len())

	l.Delete(time.Unix(0, 0))
	suite.Require().Equal(0, l.Len())
}

type loopListTestObject struct {
	Time        time.Time
	Candlestick Candlestick
}

func (suite *CandlestickListSuite) TestLoop() {
	p := "BTC-USDC"
	csList := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   p,
		Period:       period.M1,
	})

	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0))
	cs60 := Candlestick{Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(time.Unix(60, 0), cs60))
	cs120 := Candlestick{Open: 100, High: 200, Low: 50, Close: 150}
	suite.Require().NoError(csList.Set(time.Unix(120, 0), cs120))

	inspectionList := []loopListTestObject{}
	suite.Require().NoError(csList.Loop(func(ts time.Time, cs Candlestick) (bool, error) {
		inspectionList = append(inspectionList, loopListTestObject{
			Time:        ts,
			Candlestick: cs,
		})
		return false, nil
	}))

	suite.Require().Len(inspectionList, 3)
	suite.Require().Equal(inspectionList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(inspectionList[0].Candlestick, cs0)
	suite.Require().Equal(inspectionList[1].Time, time.Unix(60, 0))
	suite.Require().Equal(inspectionList[1].Candlestick, cs60)
	suite.Require().Equal(inspectionList[2].Time, time.Unix(120, 0))
	suite.Require().Equal(inspectionList[2].Candlestick, cs120)
}

func (suite *CandlestickListSuite) TestLoopError() {
	p := "BTC-USDC"
	csList := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   p,
		Period:       period.M1,
	})

	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0))
	cs60 := Candlestick{Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(time.Unix(60, 0), cs60))
	cs120 := Candlestick{Open: 100, High: 200, Low: 50, Close: 150}
	suite.Require().NoError(csList.Set(time.Unix(120, 0), cs120))

	inspectionList := []loopListTestObject{}
	suite.Require().Error(csList.Loop(func(ts time.Time, cs Candlestick) (bool, error) {
		inspectionList = append(inspectionList, loopListTestObject{
			Time:        ts,
			Candlestick: cs,
		})
		return true, errors.New("test-error")
	}))

	suite.Require().Len(inspectionList, 1)
	suite.Require().Equal(inspectionList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(inspectionList[0].Candlestick, cs0)
}

func (suite *CandlestickListSuite) TestLoopBreak() {
	p := "BTC-USDC"
	csList := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   p,
		Period:       period.M1,
	})

	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0))
	cs60 := Candlestick{Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(time.Unix(60, 0), cs60))
	cs120 := Candlestick{Open: 100, High: 200, Low: 50, Close: 150}
	suite.Require().NoError(csList.Set(time.Unix(120, 0), cs120))

	inspectionList := []loopListTestObject{}
	suite.Require().NoError(csList.Loop(func(ts time.Time, cs Candlestick) (bool, error) {
		inspectionList = append(inspectionList, loopListTestObject{
			Time:        ts,
			Candlestick: cs,
		})
		return true, nil
	}))

	suite.Require().Len(inspectionList, 1)
	suite.Require().Equal(inspectionList[0].Time, time.Unix(0, 0))
	suite.Require().Equal(inspectionList[0].Candlestick, cs0)
}

func (suite *CandlestickListSuite) TestUpdate() {
	p := "BTC-USDC"
	csList := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   p,
		Period:       period.M1,
	})

	cs0 := Candlestick{Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0))
	cs0bis := Candlestick{Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(time.Unix(0, 0), cs0bis))

	suite.Require().Equal(1, csList.Len())
	rcs0, exists := csList.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs0bis, rcs0)
}

func (suite *CandlestickListSuite) TestFirst() {
	// TODO
}

func (suite *CandlestickListSuite) TestLast() {
	// TODO
}

func (suite *CandlestickListSuite) TestReplaceUncomplete() {
	// TODO
}

func (suite *CandlestickListSuite) TestHasUncomplete() {
	// TODO
}

func (suite *CandlestickListSuite) TestExtract() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
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

func (suite *CandlestickListSuite) TestFirstN() {
	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
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

	nl := l.FirstN(2)
	suite.Require().Equal(2, nl.Len())

	cs, exists := nl.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(0.0, cs.Open)

	cs, exists = nl.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(1.0, cs.Open)
}

func (suite *CandlestickListSuite) TestMergeIntoOne() {
	// TODO
}

func (suite *CandlestickListSuite) TestLoadFromProtoBuf() {
	originalCs1 := Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}
	originalCs2 := Candlestick{
		Open:   2,
		Low:    1,
		High:   4,
		Close:  3,
		Volume: 2000,
	}

	pbList := []*proto.Candlestick{
		{
			Time:   "1970-01-01T00:01:00Z",
			Open:   1,
			Low:    0.5,
			High:   2,
			Close:  1.5,
			Volume: 1000,
		},
		{
			Time:   "1970-01-01T00:02:00Z",
			Open:   2,
			Low:    1,
			High:   4,
			Close:  3,
			Volume: 2000,
		},
	}

	l := NewEmptyList(ListID{
		ExchangeName: "exchange",
		PairSymbol:   "BTC-USDC",
		Period:       period.M1,
	})
	suite.Require().NoError(l.LoadFromProtoBuf(pbList))

	cs1, ok := l.Get(time.Unix(60, 0))
	suite.Require().True(ok)
	suite.Require().True(originalCs1.Equal(cs1))

	cs2, ok := l.Get(time.Unix(120, 0))
	suite.Require().True(ok)
	suite.Require().True(originalCs2.Equal(cs2))
}
