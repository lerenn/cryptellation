package cache

import (
	"context"
	"testing"
	"time"

	client "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

func TestCachedClient(t *testing.T) {
	suite.Run(t, new(CachedClientSuite))
}

type CachedClientSuite struct {
	exchanges    *client.MockClient
	cachedClient *cache
	suite.Suite
}

func (suite *CachedClientSuite) SetupTest() {
	suite.exchanges = client.NewMockClient(gomock.NewController(suite.T()))
	suite.cachedClient = New(suite.exchanges)
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
