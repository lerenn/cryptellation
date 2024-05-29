package client

import (
	"context"
	"testing"
	"time"

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
	candlesticks *MockClient
	cachedClient *CachedClient
	suite.Suite
}

func (suite *CachedClientSuite) SetupTest() {
	suite.candlesticks = NewMockClient(gomock.NewController(suite.T()))
	suite.cachedClient = NewCachedClient(suite.candlesticks, DefaultCacheParameters())
}

func (suite *CachedClientSuite) TestRead() {
	start := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))
	expectedRequestedStart := start.Add(-period.M1.Duration() * DefaultPreLoadingBeforeSize)
	end := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z"))
	expectedRequestedEnd := end.Add(period.M1.Duration() * DefaultPreLoadingAfterSize)

	// Setting candlesticks mock expectations
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error) {
			suite.Require().Equal("binance", payload.Exchange)
			suite.Require().Equal("BTC-USDT", payload.Pair)
			suite.Require().Equal(period.M1, payload.Period)
			suite.Require().Equal(uint(50+DefaultPreLoadingAfterSize+DefaultPreLoadingBeforeSize), payload.Limit)
			suite.Require().WithinDuration(expectedRequestedStart, *payload.Start, time.Second)
			suite.Require().WithinDuration(expectedRequestedEnd, *payload.End, time.Second)

			cl := candlestick.NewList("binance", "BTC-USDT", period.M1)
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2020-12-31T23:59:00Z")),
				candlestick.Candlestick{Close: 0}))
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")),
				candlestick.Candlestick{Close: 1}))
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:01:00Z")),
				candlestick.Candlestick{Close: 2}))
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:02:00Z")),
				candlestick.Candlestick{Close: 3}))
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z")),
				candlestick.Candlestick{Close: 4}))
			suite.Require().NoError(cl.Set(
				utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:05:00Z")),
				candlestick.Candlestick{Close: 5}))

			return cl, nil
		})

	for i := 0; i < 2; i++ {
		// Call the cached client
		cl, err := suite.cachedClient.Read(context.Background(), ReadCandlesticksPayload{
			Exchange: "binance",
			Pair:     "BTC-USDT",
			Period:   period.M1,
			Start:    utils.ToReference(start.Add(time.Nanosecond)), // Check nanosecond is not a problem
			End:      utils.ToReference(end.Add(time.Nanosecond)),   // Check nanosecond is not a problem
			Limit:    50,                                            // Check limit is not a problem
		})
		suite.Require().NoError(err)

		// Check the result
		suite.Require().Equal(4, cl.Len())

		cs, ok := cl.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(1.0, cs.Close)

		cs, ok = cl.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:01:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(2.0, cs.Close)

		cs, ok = cl.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:02:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(3.0, cs.Close)

		cs, ok = cl.Get(utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:03:00Z")))
		suite.Require().True(ok)
		suite.Require().Equal(4.0, cs.Close)
	}
}
