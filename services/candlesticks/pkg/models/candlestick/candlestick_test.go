package candlestick

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/clients/go/proto"
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

func (suite *CandlestickSuite) TestCandlestickFromProtoBuf() {
	originalCs := Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}

	t, cs, err := FromProtoBuf(&proto.Candlestick{
		Time:   "1970-01-01T00:01:00Z",
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	})

	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(60, 0), t, time.Second)
	suite.Require().True(originalCs.Equal(cs))
}

func (suite *CandlestickSuite) TestCandlestickToProtoBuf() {
	cs := Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}

	pb := cs.ToProfoBuff(time.Unix(60, 0))
	suite.Require().Equal("1970-01-01T00:01:00Z", pb.Time)
	suite.Require().Equal(float64(1), pb.Open)
	suite.Require().Equal(float64(0.5), pb.Low)
	suite.Require().Equal(float64(2), pb.High)
	suite.Require().Equal(float64(1.5), pb.Close)
	suite.Require().Equal(float64(1000), pb.Volume)
}
