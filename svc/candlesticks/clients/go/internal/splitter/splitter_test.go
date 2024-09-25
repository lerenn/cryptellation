package splitter

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestSplitterSuite(t *testing.T) {
	suite.Run(t, new(SplitterSuite))
}

type SplitterSuite struct {
	candlesticks *client.MockClient
	splitter     client.Client
	suite.Suite
}

func (suite *SplitterSuite) SetupTest() {
	suite.candlesticks = client.NewMockClient(gomock.NewController(suite.T()))
	suite.splitter = New(suite.candlesticks)
}

func (suite *SplitterSuite) TestOneCandlestick() {
	// Setting the period mock expectations
	expectedRequestedStart := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))
	expectedRequestedEnd := utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z"))

	// Set mock expectations
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			suite.Require().Equal("binance", payload.Exchange)
			suite.Require().Equal("BTC-USDT", payload.Pair)
			suite.Require().Equal(period.M1, payload.Period)
			suite.Require().Equal(uint(0), payload.Limit)
			suite.Require().WithinDuration(expectedRequestedStart, *payload.Start, time.Second)
			suite.Require().WithinDuration(expectedRequestedEnd, *payload.End, time.Second)

			cl := candlestick.NewList("binance", "BTC-USDT", period.M1)
			suite.Require().NoError(cl.Set(candlestick.Candlestick{
				Time:  utils.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")),
				Close: 1,
			}))

			return cl, nil
		})

	// Call the client
	candlesticks, err := suite.splitter.Read(context.Background(), client.ReadCandlesticksPayload{
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Period:   period.M1,
		Start:    &expectedRequestedStart,
		End:      &expectedRequestedEnd,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(1, candlesticks.Len())
}
