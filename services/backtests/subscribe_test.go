package backtests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lerenn/cryptellation/clients/go/mock"
	"github.com/lerenn/cryptellation/pkg/backtest"
	"github.com/lerenn/cryptellation/pkg/event"
	"github.com/lerenn/cryptellation/services/backtests/io/db"
	"github.com/lerenn/cryptellation/services/backtests/io/events"
	"github.com/stretchr/testify/suite"
)

func TestSubscribeSuite(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}

type SubscribeSuite struct {
	suite.Suite
	operator     Interface
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *mock.MockCandlesticks
}

func (suite *SubscribeSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *SubscribeSuite) TestHappyPath() {
	ctx := context.Background()

	// Set DB expected operations
	suite.db.EXPECT().LockedBacktest(ctx, uint(1234), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
			bt := backtest.Backtest{
				ID:                1234,
				TickSubscriptions: make([]event.TickSubscription, 0),
			}

			if err := fn(&bt); err != nil {
				return err
			}

			suite.Require().Equal(backtest.Backtest{
				ID: 1234,
				TickSubscriptions: []event.TickSubscription{
					{
						ExchangeName: "exchange",
						PairSymbol:   "ETH-USDT",
					},
				},
			}, bt)

			return nil
		})

	// Execute operation
	err := suite.operator.SubscribeToEvents(ctx, uint(1234), "exchange", "ETH-USDT")
	suite.Require().NoError(err)
}
