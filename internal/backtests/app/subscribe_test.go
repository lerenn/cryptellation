package app

import (
	"context"
	"testing"

	"github.com/digital-feather/cryptellation/clients/go/mock"
	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/events"
	"github.com/digital-feather/cryptellation/pkg/types/event"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestSubscribeSuite(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}

type SubscribeSuite struct {
	suite.Suite
	operator     Controller
	db           *db.MockAdapter
	Events       *events.MockAdapter
	candlesticks *mock.MockCandlesticks
}

func (suite *SubscribeSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.Events = events.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *SubscribeSuite) TestHappyPath() {
	ctx := context.Background()

	// Set DB expected operations
	suite.db.EXPECT().LockedBacktest(ctx, uint(1234), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
			bt := domain.Backtest{
				ID:                1234,
				TickSubscriptions: make([]event.TickSubscription, 0),
			}

			if err := fn(&bt); err != nil {
				return err
			}

			suite.Require().Equal(domain.Backtest{
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
