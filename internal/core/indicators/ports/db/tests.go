package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/stretchr/testify/suite"
)

type IndicatorsSuite struct {
	suite.Suite
	DB Port
}

func (suite *IndicatorsSuite) TestGet() {
	exchangeName := "exchange"
	pairSymbol := "ETC-USDT"
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
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	err := suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Write multiple deviating data
	p := writePayload
	p.ExchangeName = "otherExchange"
	suite.Require().NoError(suite.DB.UpsertSMA(context.Background(), p))
	p = writePayload
	p.PairSymbol = "BTC-USDC"
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
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
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
	exchangeName := "exchange"
	pairSymbol := "ETC-USDT"
	per := period.M1
	periodNumber := uint(3)
	priceType := candlestick.PriceTypeIsClose
	ts := timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2).
		Set(time.Unix(120, 0), 3).
		Set(time.Unix(180, 0), 4)
	ts2 := timeserie.New[float64]().
		Set(time.Unix(120, 0), 4).
		Set(time.Unix(180, 0), 5).
		Set(time.Unix(240, 0), 6).
		Set(time.Unix(300, 0), 7)

	// Write data from ts1
	writePayload := WriteSMAPayload{
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	err := suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Write data from ts2
	writePayload = WriteSMAPayload{
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts2,
	}
	err = suite.DB.UpsertSMA(context.Background(), writePayload)
	suite.Require().NoError(err)

	// Read data
	rts, err := suite.DB.GetSMA(context.Background(), ReadSMAPayload{
		ExchangeName: exchangeName,
		PairSymbol:   pairSymbol,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		Start:        time.Unix(0, 0),
		End:          time.Unix(180, 0),
	})
	suite.Require().NoError(err)

	// Check that data is the same
	suite.Require().Equal(ts2.Len(), rts.Len())

	// Check values
	for i := int64(0); i < 4; i++ {
		expectedValue, exists := ts2.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		value, exists := rts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		suite.Require().Equal(expectedValue, value, i)
	}
}
