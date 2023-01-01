package backtests

import (
	"context"
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
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
	candlesticks *client.MockInterfacer
}

func (suite *SubscribeSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = client.NewMockInterfacer(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
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
