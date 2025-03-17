package test

import (
	"context"
	"time"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

func (suite *EndToEndSuite) TestListIndicators() {
	// WHEN requesting for SMA

	ts, err := suite.client.ListSMA(context.Background(), api.ListSMAWorkflowParams{
		Exchange:     "binance",
		Pair:         "ETH-USDT",
		Period:       period.M1,
		Start:        utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")),
		End:          utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z")),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	})

	// THEN there is no error

	suite.Require().NoError(err)

	// AND the response contains the proper SMA

	suite.Require().Equal(3, ts.Data.Len())
	v, exists := ts.Data.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1603.8966666666668, v)
	v, exists = ts.Data.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:01:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1604.17, v)
	v, exists = ts.Data.Get(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z")))
	suite.Require().True(exists)
	suite.Require().Equal(1604.3533333333335, v)
}
