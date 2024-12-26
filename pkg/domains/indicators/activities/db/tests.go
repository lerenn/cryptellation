package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
	"github.com/stretchr/testify/suite"
)

// IndicatorsSuite is the test suite for the indicators database.
type IndicatorsSuite struct {
	suite.Suite
	DB DB
}

// TestGetSMAActivity tests the GetSMAActivity activity.
// TODO: Refactor this function
//
//nolint:funlen
func (suite *IndicatorsSuite) TestGetSMAActivity() {
	exchange := "exchange"
	pair := "ETC-USDT"
	per := period.M1
	periodNumber := 3
	priceType := candlestick.PriceTypeIsClose
	ts := timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2).
		Set(time.Unix(120, 0), 3).
		Set(time.Unix(180, 0), 4)

	// Write data
	writeParams := UpsertSMAActivityParams{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	_, err := suite.DB.UpsertSMAActivity(context.Background(), writeParams)
	suite.Require().NoError(err)

	// Write multiple deviating data
	p := writeParams
	p.Exchange = "otherExchange"
	_, err = suite.DB.UpsertSMAActivity(context.Background(), p)
	suite.Require().NoError(err)
	p = writeParams
	p.Pair = "BTC-USDC"
	_, err = suite.DB.UpsertSMAActivity(context.Background(), p)
	suite.Require().NoError(err)
	p = writeParams
	p.Period = period.D1
	_, err = suite.DB.UpsertSMAActivity(context.Background(), p)
	suite.Require().NoError(err)
	p = writeParams
	p.PeriodNumber = 8
	_, err = suite.DB.UpsertSMAActivity(context.Background(), p)
	suite.Require().NoError(err)
	p = writeParams
	p.PriceType = candlestick.PriceTypeIsOpen
	_, err = suite.DB.UpsertSMAActivity(context.Background(), p)
	suite.Require().NoError(err)

	// Read data
	rts, err := suite.DB.ReadSMAActivity(context.Background(), ReadSMAActivityParams{
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
	suite.Require().Equal(ts.Len(), rts.Data.Len())

	// Check values
	for i := int64(0); i < 4; i++ {
		expectedValue, exists := ts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		value, exists := rts.Data.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		suite.Require().Equal(expectedValue, value, i)
	}
}

// TestUpsertSMAActivity tests the UpsertSMAActivity activity.
// TODO: Refactor this function
//
//nolint:funlen
func (suite *IndicatorsSuite) TestUpsertSMAActivity() {
	exchange := "exchange"
	pair := "ETC-USDT"
	per := period.M1
	periodNumber := 3
	priceType := candlestick.PriceTypeIsClose
	ts := timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2).
		Set(time.Unix(120, 0), 3).
		Set(time.Unix(180, 0), 4)

	// Write data from ts1
	writeParams := UpsertSMAActivityParams{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	_, err := suite.DB.UpsertSMAActivity(context.Background(), writeParams)
	suite.Require().NoError(err)

	// Update data
	ts.Set(time.Unix(120, 0), 4).
		Set(time.Unix(180, 0), 5).
		Set(time.Unix(240, 0), 6).
		Set(time.Unix(300, 0), 7)

	// Write update data from ts
	writeParams = UpsertSMAActivityParams{
		Exchange:     exchange,
		Pair:         pair,
		Period:       per,
		PeriodNumber: periodNumber,
		PriceType:    priceType,
		TimeSerie:    ts,
	}
	_, err = suite.DB.UpsertSMAActivity(context.Background(), writeParams)
	suite.Require().NoError(err)

	// Read data
	rts, err := suite.DB.ReadSMAActivity(context.Background(), ReadSMAActivityParams{
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
	suite.Require().Equal(4, rts.Data.Len())

	// Check values
	for i := int64(0); i < 4; i++ {
		expectedValue, exists := ts.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		value, exists := rts.Data.Get(time.Unix(i*60, 0))
		suite.Require().True(exists, i)
		suite.Require().Equal(expectedValue, value, i)
	}
}
