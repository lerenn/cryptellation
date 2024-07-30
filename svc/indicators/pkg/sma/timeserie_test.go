package sma

import (
	"testing"
	"time"

	"cryptellation/pkg/models/timeserie"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"github.com/stretchr/testify/suite"
)

func TestTimeSerieSuite(t *testing.T) {
	suite.Run(t, new(TimeSerieSuite))
}

type TimeSerieSuite struct {
	suite.Suite
}

func (suite *TimeSerieSuite) TestTimeSerieSuite() {
	cases := []struct {
		Payload        TimeSeriePayload
		ExpectedOutput *timeserie.TimeSerie[float64]
	}{
		// Normal calculation
		{
			Payload: TimeSeriePayload{
				Candlesticks: candlestick.NewList("exchange", "ETH-USDC", period.M1).
					MustSet(time.Unix(0, 0), candlestick.Candlestick{Close: 1000}).
					MustSet(time.Unix(60, 0), candlestick.Candlestick{Close: 1500}).
					MustSet(time.Unix(120, 0), candlestick.Candlestick{Close: 1250}).
					MustSet(time.Unix(180, 0), candlestick.Candlestick{Close: 1300}),
				PriceType:    candlestick.PriceTypeIsClose,
				Start:        time.Unix(120, 0),
				End:          time.Unix(180, 0),
				PeriodNumber: 3,
			},
			ExpectedOutput: timeserie.New[float64]().
				Set(time.Unix(120, 0), 1250).
				Set(time.Unix(180, 0), 1350),
		},
	}

	for i, c := range cases {
		result := TimeSerie(c.Payload)
		suite.Require().Equal(c.ExpectedOutput.Len(), result.Len(), i)
		_ = c.ExpectedOutput.Loop(func(t time.Time, v float64) (bool, error) {
			gt, exists := result.Get(t)
			suite.Require().True(exists, i)
			suite.Require().Equal(v, gt, i)
			return false, nil
		})
	}
}

func (suite *TimeSerieSuite) TestInvalidValues() {
	cases := []struct {
		Payload        *timeserie.TimeSerie[float64]
		ExpectedOutput bool
	}{
		// No invalid values
		{
			Payload: timeserie.New[float64]().
				Set(time.Unix(0, 0), 1).
				Set(time.Unix(60, 0), 2).
				Set(time.Unix(120, 0), 3),
			ExpectedOutput: false,
		},
		// Invalid values
		{
			Payload: timeserie.New[float64]().
				Set(time.Unix(0, 0), 1).
				Set(time.Unix(60, 0), 0).
				Set(time.Unix(120, 0), 3),
			ExpectedOutput: true,
		},
	}

	for i, c := range cases {
		suite.Require().Equal(c.ExpectedOutput, InvalidValues(c.Payload), i)
	}
}
