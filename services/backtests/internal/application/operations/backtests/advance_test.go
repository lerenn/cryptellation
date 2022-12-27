package backtests

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestAdvanceSuite(t *testing.T) {
	suite.Run(t, new(AdvanceSuite))
}

type AdvanceSuite struct {
	suite.Suite
	operator     Operator
	db           *db.MockAdapter
	pubsub       *pubsub.MockAdapter
	candlesticks *client.MockInterfacer
}

func (suite *AdvanceSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = client.NewMockInterfacer(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
}

func (suite *AdvanceSuite) TestWithoutAccount() {
	ctx := context.Background()

	// Set DB calls expectated
	suite.db.EXPECT().LockedBacktest(uint(1234), gomock.Any()).
		DoAndReturn(func(id uint, fn db.LockedBacktestCallback) error {
			return fn()
		})

	suite.db.EXPECT().ReadBacktest(ctx, uint(1234)).Return(backtest.Backtest{
		ID: 1234,
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(0, 0),
			PriceType: candlestick.PriceTypeIsOpen,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		TickSubscribers: []event.Subscription{
			{
				ExchangeName: "exchange",
				PairSymbol:   "ETH-USDT",
			},
		},
	}, nil)

	suite.db.EXPECT().UpdateBacktest(gomock.Any(), backtest.Backtest{
		ID: 1234,
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(120, 0),
			PriceType: candlestick.PriceTypeIsOpen,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		TickSubscribers: []event.Subscription{
			{
				ExchangeName: "exchange",
				PairSymbol:   "ETH-USDT",
			},
		},
	})

	// Set candlesticks client expected calls
	suite.candlesticks.EXPECT().ReadCandlesticks(ctx, client.ReadCandlestickPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDT",
		Period:       period.M1,
		Start:        time.Unix(0, 0),
		End:          time.Unix(120, 0),
		Limit:        1,
	}).Return(candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDT",
		Period:       period.M1,
	}), nil)

	// Set pubsub expected calls
	suite.pubsub.EXPECT().Publish(uint(1234), event.NewStatusEvent(time.Unix(120, 0), status.Status{Finished: true}))

	// Execute operation
	err := suite.operator.Advance(context.Background(), uint(1234))

	// Check return
	suite.Require().NoError(err)
}

func (suite *AdvanceSuite) TestWithAnAccount() {
	// TODO
}

func (suite *AdvanceSuite) TestWithAnOrder() {
	// TODO
}
