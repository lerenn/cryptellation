package app

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/clients/go/mock"
	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/events"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/pkg/types/event"
	"github.com/digital-feather/cryptellation/pkg/types/order"
	"github.com/digital-feather/cryptellation/pkg/types/period"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestCreationSuite(t *testing.T) {
	suite.Run(t, new(CreationSuite))
}

type CreationSuite struct {
	suite.Suite
	operator     Controller
	db           *db.MockAdapter
	Events       *events.MockAdapter
	candlesticks *mock.MockCandlesticks
}

func (suite *CreationSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.Events = events.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *CreationSuite) TestHappyPass() {
	ctx := context.Background()

	// Set DB mock expectations
	suite.db.EXPECT().CreateBacktest(ctx, gomock.Eq(&domain.Backtest{
		StartTime: time.Unix(0, 0),
		CurrentCsTick: domain.CurrentCsTick{
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
		Do(func(ctx context.Context, bt *domain.Backtest) { bt.ID = 1 }).
		Return(nil)

	// Execute creation
	id, err := suite.operator.Create(ctx, domain.NewPayload{
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
