package domain

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestSubscribeSuite(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}

type SubscribeSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *SubscribeSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
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
						Exchange: "exchange",
						Pair:     "ETH-USDT",
					},
				},
			}, bt)

			return nil
		})

	// Execute operation
	err := suite.operator.SubscribeToEvents(ctx, uint(1234), "exchange", "ETH-USDT")
	suite.Require().NoError(err)
}
