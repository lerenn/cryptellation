package client

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

func TestCachedClientSuite(t *testing.T) {
	suite.Run(t, new(CachedClientSuite))
}

type CachedClientSuite struct {
	indicators   *MockClient
	cachedClient *CachedClient
	suite.Suite
}

func (suite *CachedClientSuite) SetupTest() {
	suite.indicators = NewMockClient(gomock.NewController(suite.T()))
	suite.cachedClient = NewCachedClient(suite.indicators, DefaultCacheParameters())
}

func (suite *CachedClientSuite) TeardownTest() {
	suite.indicators.EXPECT().Close(gomock.Any()).Return()
	suite.cachedClient.Close(context.Background())
}

func (suite *CachedClientSuite) TestSMA() {
	start := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))
	expectedRequestedStart := start.Add(-period.M1.Duration() * DefaultPreLoadingBeforeSize)
	end := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z"))
	expectedRequestedEnd := end.Add(period.M1.Duration() * DefaultPreLoadingAfterSize)

	// Setting indicators mock expectations
	suite.indicators.EXPECT().SMA(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload SMAPayload) (*timeserie.TimeSerie[float64], error) {
			suite.Require().Equal("binance", payload.Exchange)
			suite.Require().Equal("BTC-USDT", payload.Pair)
			suite.Require().Equal(period.M1, payload.Period)
			suite.Require().Equal(uint(14), payload.PeriodNumber)
			suite.Require().Equal(candlestick.PriceTypeIsClose, payload.PriceType)
			suite.Require().WithinDuration(expectedRequestedStart, payload.Start, time.Second)
			suite.Require().WithinDuration(expectedRequestedEnd, payload.End, time.Second)

			ts := timeserie.New[float64]()
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2020-12-31T23:59:00Z")), 0.0)
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")), 1.0)
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:01:00Z")), 2.0)
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:02:00Z")), 3.0)
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z")), 4.0)
			_ = ts.Set(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:04:00Z")), 5.0)
			return ts, nil
		})

	// Testing SMA twice
	for i := 0; i < 2; i++ {
		ts, err := suite.cachedClient.SMA(context.Background(), SMAPayload{
			Exchange:     "binance",
			Pair:         "BTC-USDT",
			Period:       period.M1,
			Start:        start.Add(time.Nanosecond), // Check nanosecond is not a problem
			End:          end.Add(time.Nanosecond),   // Check nanosecond is not a problem
			PeriodNumber: 14,
			PriceType:    candlestick.PriceTypeIsClose,
		})
		suite.Require().NoError(err)
		suite.Require().Equal(4, ts.Len())

		p, ok := ts.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(1.0, p)

		p, ok = ts.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:01:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(2.0, p)

		p, ok = ts.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:02:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(3.0, p)

		p, ok = ts.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(4.0, p)
	}
}
