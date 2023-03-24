package backtests

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/clients/go/mock"
	"github.com/digital-feather/cryptellation/pkg/account"
	"github.com/digital-feather/cryptellation/pkg/backtest"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/event"
	"github.com/digital-feather/cryptellation/pkg/order"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/services/backtests/io/db"
	"github.com/digital-feather/cryptellation/services/backtests/io/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestCreationSuite(t *testing.T) {
	suite.Run(t, new(CreationSuite))
}

type CreationSuite struct {
	suite.Suite
	operator     Interface
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *mock.MockCandlesticks
}

func (suite *CreationSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *CreationSuite) TestHappyPass() {
	ctx := context.Background()

	// Set DB mock expectations
	suite.db.EXPECT().CreateBacktest(ctx, gomock.Eq(&backtest.Backtest{
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(0, 0),
			PriceType: candlestick.PriceTypeIsOpen,
		},
		EndTime: time.Unix(120, 0),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{"DAI": 1000},
			},
		},
		PeriodBetweenEvents: period.M1,
		TickSubscriptions:   make([]event.TickSubscription, 0),
		Orders:              make([]order.Order, 0)})).
		Do(func(ctx context.Context, bt *backtest.Backtest) { bt.ID = 1 }).
		Return(nil)

	// Execute creation
	id, err := suite.operator.Create(ctx, backtest.NewPayload{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{"DAI": 1000},
			},
		},
		StartTime:             time.Unix(0, 0),
		EndTime:               TimeOpt(time.Unix(120, 0)),
		DurationBetweenEvents: DurationOpt(time.Minute),
	})

	// Check that returned value is correct
	suite.Require().Equal(uint(1), id)
	suite.Require().NoError(err)
}

func TimeOpt(t time.Time) *time.Time {
	return &t
}

func DurationOpt(t time.Duration) *time.Duration {
	return &t
}
