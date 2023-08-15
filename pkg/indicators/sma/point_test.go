package sma

import (
	"math"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/stretchr/testify/suite"
)

func TestPointSuite(t *testing.T) {
	suite.Run(t, new(PointSuite))
}

type PointSuite struct {
	suite.Suite
}

func (suite *PointSuite) TestPoint() {
	cases := []struct {
		Payload        PointPayload
		ExpectedOutput float64
	}{
		// Normal calculation
		{
			Payload: PointPayload{
				Candlesticks: timeserie.New[candlestick.Candlestick]().
					Set(time.Unix(0, 0), candlestick.Candlestick{Close: 1000}).
					Set(time.Unix(60, 0), candlestick.Candlestick{Close: 1500}).
					Set(time.Unix(120, 0), candlestick.Candlestick{Close: 1250}),
				PriceType: candlestick.PriceTypeIsClose,
			},
			ExpectedOutput: 1250,
		},
		// No point
		{
			Payload: PointPayload{
				Candlesticks: timeserie.New[candlestick.Candlestick](),
				PriceType:    candlestick.PriceTypeIsClose,
			},
			ExpectedOutput: math.NaN(),
		},
	}

	for i, c := range cases {
		if math.IsNaN(c.ExpectedOutput) {
			suite.Require().True(math.IsNaN(Point(c.Payload)), i)
		} else {
			suite.Require().Equal(c.ExpectedOutput, Point(c.Payload), i)
		}
	}
}
