package backtests

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestCreationSuite(t *testing.T) {
	suite.Run(t, new(CreationSuite))
}

type CreationSuite struct {
	suite.Suite
	operator     Operator
	db           *db.MockAdapter
	pubsub       *pubsub.MockAdapter
	candlesticks *client.MockInterfacer
}

func (suite *CreationSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = client.NewMockInterfacer(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
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
		TickSubscribers:     make([]event.Subscription, 0),
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
