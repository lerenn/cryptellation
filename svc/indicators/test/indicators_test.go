package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	client "github.com/lerenn/cryptellation/svc/indicators/clients/go"
)

func (suite *EndToEndSuite) TestGetSMA() {
	// WHEN requesting for SMA
	ts, err := suite.client.SMA(context.Background(), client.SMAPayload{
		Exchange:     "binance",
		Pair:         "ETH-USDT",
		Period:       period.M1,
		Start:        utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")),
		End:          utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z")),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	})

	// THEN response is adequate
	suite.Require().NoError(err)
	suite.Require().Equal(3, ts.Len())
	v, exists := ts.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1603.8966666666668, v)
	v, exists = ts.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:01:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1604.17, v)
	v, exists = ts.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1604.3533333333335, v)
}
