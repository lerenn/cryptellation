package domain

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/models/event"

	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"

	"github.com/google/uuid"
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
	id := uuid.New()

	// Set DB expected operations
	suite.db.EXPECT().LockedBacktest(ctx, id, gomock.Any()).
		DoAndReturn(func(ctx context.Context, id uuid.UUID, fn db.LockedBacktestCallback) error {
			bt := backtest.Backtest{
				ID:                uuid.New(),
				TickSubscriptions: make([]event.TickSubscription, 0),
			}

			if err := fn(&bt); err != nil {
				return err
			}

			suite.Require().Equal(backtest.Backtest{
				ID: bt.ID,
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
	err := suite.operator.SubscribeToEvents(ctx, id, "exchange", "ETH-USDT")
	suite.Require().NoError(err)
}
