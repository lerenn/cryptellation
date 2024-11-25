package candlestick

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/period"
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
	suite.Require().Equal("exchange", l.Metadata.Exchange)
	suite.Require().Equal("ETH-USDC", l.Metadata.Pair)
	suite.Require().Equal(period.M1, l.Metadata.Period)
	suite.Require().Equal(0, l.Data.Len())
}

func (suite *CandlestickListSuite) TestNew() {
	cs1 := Candlestick{
		Time: time.Unix(0, 0),
		Open: 1.0,
	}
	cs2 := Candlestick{
		Time: time.Unix(60, 0),
		Open: 2.0,
	}
	l := NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(l.Set(cs1))
	suite.Require().NoError(l.Set(cs2))

	// Check list
	suite.Require().Equal("exchange", l.Metadata.Exchange)
	suite.Require().Equal("ETH-USDC", l.Metadata.Pair)
	suite.Require().Equal(period.M1, l.Metadata.Period)

	// Check candlesticks
	suite.Require().Equal(2, l.Data.Len())

	t, e := l.Data.Get(cs1.Time)
	suite.Require().True(e)
	suite.Require().Equal(cs1, t)

	t, e = l.Data.Get(cs2.Time)
	suite.Require().True(e)
	suite.Require().Equal(cs2, t)
}

func (suite *CandlestickListSuite) TestJSONMarshaling() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Time: time.Unix(0, 0),
		Open: 1.0,
	}
	suite.Require().NoError(l.Set(cs))

	// Marshal
	data, err := json.Marshal(l)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(data)

	// Unmarshal
	l2 := &List{}
	err = json.Unmarshal(data, &l2)
	suite.Require().NoError(err)
	suite.Require().Equal(l.Metadata, l2.Metadata)
	t1, _ := l.Data.Get(time.Unix(0, 0))
	t2, _ := l.Data.Get(time.Unix(0, 0))
	suite.Require().Equal(t1, t2)
}

func (suite *CandlestickListSuite) TestNewWithUnalignedCandlestick() {
	cs := Candlestick{
		Time: time.Unix(1, 0),
		Open: 1.0,
	}
	l := NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().Error(l.Set(cs))
}

func (suite *CandlestickListSuite) TestMustSet() {
	// TODO: set
}

func (suite *CandlestickListSuite) TestSet() {
	p := "BTC-USDC"
	csList := NewList("exchange", p, period.M1)

	cs0 := Candlestick{Time: time.Unix(0, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(cs0))
	cs0bis := Candlestick{Time: time.Unix(0, 0), Open: 10, High: 20, Low: 5, Close: 15}
	suite.Require().NoError(csList.Set(cs0bis))

	suite.Require().Equal(1, csList.Data.Len())
	rcs0, exists := csList.Data.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs0bis, rcs0)
}

func (suite *CandlestickListSuite) TestSetWithWrongPeriod() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Time:  time.Unix(1, 0),
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	suite.Require().Error(l.Set(cs))
}

func (suite *CandlestickListSuite) TestMerge() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	recvCSList := NewList("exchange", "ETH-USDC", period.M1)
	cs := Candlestick{
		Time:  time.Unix(0, 0),
		Open:  1,
		High:  2,
		Low:   0.5,
		Close: 1.5,
	}
	err := recvCSList.Set(cs)
	suite.Require().NoError(err)
	suite.Require().Equal(1, recvCSList.Data.Len())

	err = l.Merge(recvCSList, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(1, l.Data.Len())
	cs2, exists := l.Data.Get(time.Unix(0, 0))
	suite.Require().True(exists)
	suite.Require().Equal(cs, cs2)
}

func (suite *CandlestickListSuite) TestExtract() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	for i := int64(0); i < 4; i++ {
		cs := Candlestick{
			Time:  time.Unix(60*i, 0),
			Open:  float64(i),
			High:  0,
			Low:   0,
			Close: 0,
		}

		err := l.Set(cs)
		suite.Require().NoError(err)
	}

	nl := l.Extract(time.Unix(60, 0), time.Unix(120, 0), 0)
	suite.Require().Equal(2, nl.Data.Len())

	cs, exists := nl.Data.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(1.0, cs.Open)

	cs, exists = nl.Data.Get(time.Unix(120, 0))
	suite.Require().True(exists)
	suite.Require().Equal(2.0, cs.Open)
}

func (suite *CandlestickListSuite) TestExtractWithLimit() {
	l := NewList("exchange", "ETH-USDC", period.M1)
	for i := int64(0); i < 4; i++ {
		cs := Candlestick{
			Time:  time.Unix(60*i, 0),
			Open:  float64(i),
			High:  0,
			Low:   0,
			Close: 0,
		}

		err := l.Set(cs)
		suite.Require().NoError(err)
	}

	nl := l.Extract(time.Unix(60, 0), time.Unix(180, 0), 2)
	suite.Require().Equal(2, nl.Data.Len())

	cs, exists := nl.Data.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(1.0, cs.Open)

	cs, exists = nl.Data.Get(time.Unix(120, 0))
	suite.Require().True(exists)
	suite.Require().Equal(2.0, cs.Open)
}

func (suite *CandlestickListSuite) TestMergeWithNotCorrespondingLists() {
	l1 := NewList("exchange", "ETH-USDC", period.M1)
	l2 := NewList("exchange", "ETH-USDC", period.M5)
	err := l1.Merge(l2, nil)
	suite.Require().Error(err)
}

func (suite *CandlestickListSuite) TestReplaceUncomplete() {
	// Create a list with one candlestick
	completeList := NewList("exchange", "BTC-USDC", period.M1)
	csComplete := Candlestick{Time: time.Unix(120, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(completeList.Set(csComplete))

	// Create a list with one uncomplete candlestick
	uncompleteList := NewList("exchange", "BTC-USDC", period.M1)
	csUncomplete := Candlestick{Time: time.Unix(120, 0), Open: 10, High: 20, Low: 5, Close: 15, Uncomplete: true}
	suite.Require().NoError(uncompleteList.Set(csUncomplete))

	// Replace uncomplete candlestick
	uncompleteList.ReplaceUncomplete(completeList)

	// Check that the uncomplete candlestick has been replaced
	cs, exists := uncompleteList.Data.Get(time.Unix(120, 0))
	suite.Require().True(exists)
	suite.Require().Equal(csComplete, cs)
}

func (suite *CandlestickListSuite) TestGetUncompleteTimes() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	cs0 := Candlestick{Time: time.Unix(60, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(cs0))

	// Get uncomplete times
	times := csList.GetUncompleteTimes()
	suite.Require().Empty(times)

	// Add an uncomplete candlestick
	cs1 := Candlestick{Time: time.Unix(120, 0), Open: 10, High: 20, Low: 5, Close: 15, Uncomplete: true}
	suite.Require().NoError(csList.Set(cs1))

	// Get uncomplete times
	times = csList.GetUncompleteTimes()
	suite.Require().Len(times, 1)
	suite.Require().Equal(time.Unix(120, 0), times[0])
}

func (suite *CandlestickListSuite) TestFillMissing() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	cs0 := Candlestick{Time: time.Unix(60, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(cs0))

	// Fill missing candlesticks
	err := csList.FillMissing(time.Unix(0, 0), time.Unix(180, 0), Candlestick{Open: 10, High: 20, Low: 5, Close: 15})
	suite.Require().NoError(err)

	// Check that the existing one is not overwritten
	cs, exists := csList.Data.Get(time.Unix(60, 0))
	suite.Require().True(exists)
	suite.Require().Equal(Candlestick{Time: time.Unix(60, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}, cs)

	// Check that missing candlesticks has been filled
	for _, m := range []int64{0, 2, 3} {
		t := time.Unix(m*60, 0)
		cs, exists := csList.Data.Get(t)
		suite.Require().True(exists)
		suite.Require().Equal(Candlestick{Time: t, Open: 10, High: 20, Low: 5, Close: 15}, cs)
	}
}

func (suite *CandlestickListSuite) TestFillMissingWithEmptyList() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	err := csList.FillMissing(time.Unix(0, 0), time.Unix(180, 0), Candlestick{})
	suite.Require().NoError(err)
}

func (suite *CandlestickListSuite) TestGetMissingTimes() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	cs0 := Candlestick{Time: time.Unix(60, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(cs0))

	// Get missing times
	times := csList.GetMissingTimes(time.Unix(0, 0), time.Unix(180, 0), 0)
	suite.Require().Len(times, 3)
	suite.Require().Equal(time.Unix(0, 0), times[0])
	suite.Require().Equal(time.Unix(120, 0), times[1])
	suite.Require().Equal(time.Unix(180, 0), times[2])
}

func (suite *CandlestickListSuite) TestToArray() {
	// Create a list with one candlestick
	csList := NewList("exchange", "BTC-USDC", period.M1)
	cs0 := Candlestick{Time: time.Unix(60, 0), Open: 1, High: 2, Low: 0.5, Close: 1.5}
	suite.Require().NoError(csList.Set(cs0))
	cs1 := Candlestick{Time: time.Unix(120, 0), Open: 2, High: 4, Low: 1, Close: 3}
	suite.Require().NoError(csList.Set(cs1))

	// Convert to array
	arr := csList.ToArray()
	suite.Require().Len(arr, 2)
	suite.Require().Equal(cs0, arr[0])
	suite.Require().Equal(cs1, arr[1])
}
