package sma

import (
	"math"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

func TestPointSuite(t *testing.T) {
	suite.Run(t, new(PointSuite))
}

type PointSuite struct {
	suite.Suite
}

func (suite *PointSuite) TestPoint() {
	csNormal := candlestick.NewList("binance", "ETH-USDT", period.M1)
	suite.Require().NoError(csNormal.Set(candlestick.Candlestick{Time: time.Unix(0, 0), Close: 1000}))
	suite.Require().NoError(csNormal.Set(candlestick.Candlestick{Time: time.Unix(60, 0), Close: 1500}))
	suite.Require().NoError(csNormal.Set(candlestick.Candlestick{Time: time.Unix(120, 0), Close: 1250}))

	csMissingPrice := candlestick.NewList("binance", "ETH-USDT", period.M1)
	suite.Require().NoError(csMissingPrice.Set(candlestick.Candlestick{Time: time.Unix(0, 0), Close: 1000}))
	suite.Require().NoError(csMissingPrice.Set(candlestick.Candlestick{Time: time.Unix(60, 0), Close: 1500}))
	suite.Require().NoError(csMissingPrice.Set(candlestick.Candlestick{Time: time.Unix(120, 0), Close: 1250}))

	cases := []struct {
		Params         PointParameters
		ExpectedOutput float64
	}{
		// Normal calculation
		{
			Params: PointParameters{
				Candlesticks: csNormal,
				PriceType:    candlestick.PriceTypeIsClose,
			},
			ExpectedOutput: 1250,
		},
		// Calculation with missing price
		{
			Params: PointParameters{
				Candlesticks: csMissingPrice,
				PriceType:    candlestick.PriceTypeIsClose,
			},
			ExpectedOutput: 1125,
		},
		// No point
		{
			Params: PointParameters{
				Candlesticks: candlestick.NewList("binance", "ETH-USDT", period.M1),
				PriceType:    candlestick.PriceTypeIsClose,
			},
			ExpectedOutput: math.NaN(),
		},
	}

	for i, c := range cases {
		p, err := NewPoint(c.Params)
		suite.Require().NoError(err, i)

		if math.IsNaN(c.ExpectedOutput) {
			suite.Require().True(math.IsNaN(p.Price), i)
		} else {
			suite.Require().Equal(c.ExpectedOutput, p.Price, i)
		}
	}
}
