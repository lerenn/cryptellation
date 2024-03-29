package domain

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestCreationSuite(t *testing.T) {
	suite.Run(t, new(CreationSuite))
}

type CreationSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *CreationSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
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
