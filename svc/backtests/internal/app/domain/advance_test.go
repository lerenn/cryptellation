package domain

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestAdvanceSuite(t *testing.T) {
	suite.Run(t, new(AdvanceSuite))
}

type AdvanceSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *AdvanceSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.events, suite.candlesticks)
}

func (suite *AdvanceSuite) TestWithoutAccount() {
	ctx := context.Background()

	// Set DB calls expectated
	suite.db.EXPECT().LockedBacktest(ctx, uint(1234), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
			bt := backtest.Backtest{
				ID: 1234,
				CurrentCsTick: backtest.CurrentCsTick{
					Time:      time.Unix(0, 0),
					PriceType: candlestick.PriceTypeIsOpen,
				},
				EndTime:             time.Unix(120, 0),
				PeriodBetweenEvents: period.M1,
				TickSubscriptions: []event.TickSubscription{
					{
						Exchange: "exchange",
						Pair:     "ETH-USDT",
					},
				},
			}

			if err := fn(&bt); err != nil {
				return err
			}

			suite.Require().Equal(backtest.Backtest{
				ID: 1234,
				CurrentCsTick: backtest.CurrentCsTick{
					Time:      time.Unix(120, 0),
					PriceType: candlestick.PriceTypeIsOpen,
				},
				EndTime:             time.Unix(120, 0),
				PeriodBetweenEvents: period.M1,
				TickSubscriptions: []event.TickSubscription{
					{
						Exchange: "exchange",
						Pair:     "ETH-USDT",
					},
				},
			}, bt)

			return nil
		})

	// Set candlesticks client expected calls
	suite.candlesticks.EXPECT().Read(ctx, candlesticks.ReadCandlesticksPayload{
		Exchange: "exchange",
		Pair:     "ETH-USDT",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(0, 0)),
		End:      utils.ToReference(time.Unix(120, 0)),
		Limit:    1,
	}).Return(candlestick.NewList("exchange", "ETH-USDT", period.M1), nil)

	// Set Events expected calls
	suite.events.EXPECT().Publish(context.Background(), uint(1234), event.NewStatusEvent(time.Unix(120, 0), event.Status{Finished: true}))

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
