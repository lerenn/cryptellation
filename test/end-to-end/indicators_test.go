package endToEnd

import (
	"context"
	"testing"
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestIndicatorsSuite(t *testing.T) {
	suite.Run(t, new(IndicatorsSuite))
}

type IndicatorsSuite struct {
	suite.Suite
	client client.Indicators
}

func (suite *IndicatorsSuite) SetupSuite() {
	// Get config
	cfg := config.LoadDefaultNATSConfig()
	cfg.OverrideFromEnv()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := nats.NewIndicators(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *IndicatorsSuite) TearDownSuite() {
	suite.client.Close(context.Background())
}

func (suite *IndicatorsSuite) TestGetSMA() {
	// WHEN requesting for SMA
	ts, err := suite.client.SMA(context.Background(), client.SMAPayload{
		ExchangeName: "binance",
		PairSymbol:   "ETH-USDT",
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
