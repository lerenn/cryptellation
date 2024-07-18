package retry

import (
	"context"
	"testing"
	"time"

	common "github.com/lerenn/cryptellation/pkg/client"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

func TestRetrySuite(t *testing.T) {
	suite.Run(t, new(RetrySuite))
}

type RetrySuite struct {
	candlesticks *client.MockClient
	Retry        client.Client
	suite.Suite
}

func (suite *RetrySuite) SetupTest() {
	suite.candlesticks = client.NewMockClient(gomock.NewController(suite.T()))
	suite.Retry = New(suite.candlesticks,
		WithMaxRetry(3),
		WithTimeout(time.Millisecond*100),
	)
}

func (suite *RetrySuite) TestReadUntilExpiration() {
	// Setting candlesticks mock expectations with only
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})

	// Testing the Read method
	_, err := suite.Retry.Read(context.Background(), client.ReadCandlesticksPayload{})
	suite.ErrorIs(err, common.ErrMaxRetriesReached)
}

func (suite *RetrySuite) TestReadWithOneExpiration() {
	// Setting candlesticks mock expectations with only
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})
	suite.candlesticks.EXPECT().Read(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			return candlestick.NewList("binance", "BTC-USDT", period.M1), nil
		})

	// Testing the Read method
	_, err := suite.Retry.Read(context.Background(), client.ReadCandlesticksPayload{})
	suite.NoError(err)
}
