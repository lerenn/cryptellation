package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
)

type IndicatorsSuite struct {
	suite.Suite
	DB Port
}

func (suite *IndicatorsSuite) TestGet() {
	exchange := "exchange"
	pair := "ETC-USDT"
	per := period.M1
	periodNumber := uint(3)
	priceType := candlestick.PriceTypeIsClose
	ts := timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2).
		Set(time.Unix(120, 0), 3).
		Set(time.Unix(180, 0), 4)

	// Write data
	writePayload := WriteSMAPayload{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	err := suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Write multiple deviating data
	p := writePayload
	p.Exchange = "otherExchange"
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))
	p = writePayload
	p.Pair = "BTC-USDC"
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))
	p = writePayload
	p.Period = period.D1
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))
	p = writePayload
	p.PeriodNumber = 8
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))
	p = writePayload
	p.PriceType = candlestick.PriceTypeIsOpen
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))

	// Read data
	rts, err := suite.DB.GetSMA(context.Background(), ReadSMAPayload{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		Start:        time.Unix(0, 0),
		End:          time.Unix(180, 0),
	})
	suite.Require().NoError(err)

	// Check that data is the same
	suite.Require().Equal(ts.Len(), rts.Len())

	// Check values
	for i := int64(0); i < 4; i++ {
		expectedValue, exists := ts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		value, exists := rts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		suite.Require().Equal(expectedValue, value, i)
	}
}

func (suite *IndicatorsSuite) TestUpsert() {
	exchange := "exchange"
	pair := "ETC-USDT"
	per := period.M1
	periodNumber := uint(3)
	priceType := candlestick.PriceTypeIsClose
	ts := timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2).
		Set(time.Unix(120, 0), 3).
		Set(time.Unix(180, 0), 4)

	// Write data from ts1
	writePayload := WriteSMAPayload{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	err := suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Update data
	ts.Set(time.Unix(120, 0), 4).
		Set(time.Unix(180, 0), 5).
		Set(time.Unix(240, 0), 6).
		Set(time.Unix(300, 0), 7)

	// Write update data from ts
	writePayload = WriteSMAPayload{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	err = suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Read data
	rts, err := suite.DB.GetSMA(context.Background(), ReadSMAPayload{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		Start:        time.Unix(0, 0),
		End:          time.Unix(180, 0),
	})
	suite.Require().NoError(err)

	// Check that data is the same
	suite.Require().Equal(4, rts.Len())

	// Check values
	for i := int64(0); i < 4; i++ {
		expectedValue, exists := ts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		value, exists := rts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		suite.Require().Equal(expectedValue, value, i)
	}
}
