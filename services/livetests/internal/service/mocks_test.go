package service

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/pairs"
	ticksProto "github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
	"github.com/stretchr/testify/suite"
)

func TestMocksSuite(t *testing.T) {
	suite.Run(t, new(MocksSuite))
}

type MocksSuite struct {
	suite.Suite
}

func (suite *MocksSuite) TestTicks() {
	var client MockedTicksClient

	req := ticksProto.ListenSymbolRequest{
		Exchange:   "mocked_exchange",
		PairSymbol: pairs.FormatPairSymbol("ETH", "USDC"),
	}
	tClient, err := client.ListenSymbol(context.TODO(), &req)
	suite.Require().NoError(err)

	for i := 0; i < 100; i++ {
		t, err := tClient.Recv()
		suite.Require().NoError(err)

		suite.Require().Equal(req.Exchange, t.Exchange)
		suite.Require().Equal(req.PairSymbol, t.PairSymbol)
		suite.Require().Equal(time.Unix(int64(i), 0).Format(time.RFC3339), t.Time)
		suite.Require().Equal(float32(i), t.Price)
	}
}
