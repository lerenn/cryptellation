package client

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

func TestCachedClient(t *testing.T) {
	suite.Run(t, new(CachedClientSuite))
}

type CachedClientSuite struct {
	exchanges    *MockClient
	cachedClient *CachedClient
	suite.Suite
}

func (suite *CachedClientSuite) SetupTest() {
	suite.exchanges = NewMockClient(gomock.NewController(suite.T()))
	suite.cachedClient = NewCachedClient(suite.exchanges, DefaultCacheParameters())
}

func (suite *CachedClientSuite) TestRead() {
	// Setting exchanges mock expectations
	suite.exchanges.EXPECT().Read(context.Background(), []string{"binance", "coinbase"}).
		Return([]exchange.Exchange{
			{Name: "binance", LastSyncTime: time.Now()},
			{Name: "coinbase", LastSyncTime: time.Now()},
		}, nil)

	// Reading exchanges
	for i := 0; i < 2; i++ {
		exchanges, err := suite.cachedClient.Read(context.Background(), "binance", "coinbase")
		suite.Require().NoError(err)
		suite.Require().Len(exchanges, 2)
		suite.Require().Equal("binance", exchanges[0].Name)
		suite.Require().Equal("coinbase", exchanges[1].Name)
	}
}
