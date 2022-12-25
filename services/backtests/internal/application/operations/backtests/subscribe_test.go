package backtests

import (
	"context"
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestSubscribeSuite(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}

type SubscribeSuite struct {
	suite.Suite
	operator     Operator
	db           *db.MockAdapter
	pubsub       *pubsub.MockAdapter
	candlesticks *MockClient
}

func (suite *SubscribeSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
}

func (suite *SubscribeSuite) TestHappyPath() {
	ctx := context.Background()

	// Set DB expected operations
	suite.db.EXPECT().LockedBacktest(uint(1234), gomock.Any()).
		DoAndReturn(func(id uint, fn db.LockedBacktestCallback) error {
			return fn()
		})
	suite.db.EXPECT().ReadBacktest(ctx, uint(1234)).Return(backtest.Backtest{
		ID:              1234,
		TickSubscribers: make([]event.Subscription, 0),
	}, nil)
	suite.db.EXPECT().UpdateBacktest(ctx, backtest.Backtest{
		ID: 1234,
		TickSubscribers: []event.Subscription{
			{
				ExchangeName: "exchange",
				PairSymbol:   "ETH-USDT",
			},
		},
	})

	// Execute operation
	err := suite.operator.SubscribeToEvents(ctx, uint(1234), "exchange", "ETH-USDT")
	suite.Require().NoError(err)
}
