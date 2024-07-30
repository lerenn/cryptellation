package domain

import (
	"context"
	"testing"
	"time"

	"cryptellation/pkg/models/account"
	"cryptellation/pkg/models/order"

	candlesticks "cryptellation/svc/candlesticks/clients/go"
	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"cryptellation/svc/forwardtests/internal/app"
	"cryptellation/svc/forwardtests/internal/app/ports/db"
	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestOrderCreationSuite(t *testing.T) {
	suite.Run(t, new(OrderCreationSuite))
}

type OrderCreationSuite struct {
	suite.Suite
	operator     app.ForwardTests
	db           *db.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *OrderCreationSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.candlesticks)
}

func (suite *OrderCreationSuite) TestHappyPass() {
	// Order creation parameters
	id := uuid.New()
	o := order.Order{
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Quantity: 1,
		Type:     order.TypeIsMarket,
		Side:     order.SideIsBuy,
	}

	// Setting candlesticks mock expectations
	suite.candlesticks.EXPECT().Read(context.Background(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload candlesticks.ReadCandlesticksPayload) (*candlestick.List, error) {
			suite.Require().Equal("binance", payload.Exchange)
			suite.Require().Equal("BTC-USDT", payload.Pair)
			suite.Require().WithinDuration(time.Now(), *payload.Start, time.Second)
			suite.Require().WithinDuration(time.Now(), *payload.End, time.Second)
			suite.Require().Equal(uint(1), payload.Limit)

			cl := candlestick.NewList("binance", "BTC-USDT", period.M1)
			err := cl.Set(time.Now().Round(time.Minute), candlestick.Candlestick{
				Open:  800,
				Close: 1000,
				High:  1200,
				Low:   700,
			})
			suite.Require().NoError(err)

			return cl, nil
		})

	// Setting DB mock expectations
	suite.db.EXPECT().ReadForwardTest(context.Background(), id).
		Return(forwardtest.ForwardTest{
			ID: id,
			Accounts: map[string]account.Account{
				"binance": {
					Balances: map[string]float64{"USDT": 1000},
				},
			},
			Orders: []order.Order{},
		}, nil)
	suite.db.EXPECT().UpdateForwardTest(context.Background(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ft forwardtest.ForwardTest) error {
			suite.Require().Equal(id, ft.ID)
			suite.Require().Equal(map[string]account.Account{
				"binance": {
					Balances: map[string]float64{"USDT": 0, "BTC": 1},
				},
			}, ft.Accounts)

			suite.Require().Len(ft.Orders, 1)
			suite.Require().NotEqual(uuid.Nil, ft.Orders[0].ID)
			suite.Require().Equal("binance", ft.Orders[0].Exchange)
			suite.Require().Equal("BTC-USDT", ft.Orders[0].Pair)
			suite.Require().Equal(1.0, ft.Orders[0].Quantity)
			suite.Require().Equal(order.TypeIsMarket, ft.Orders[0].Type)
			suite.Require().Equal(order.SideIsBuy, ft.Orders[0].Side)
			suite.Require().NotNil(ft.Orders[0].ExecutionTime)
			suite.Require().WithinDuration(time.Now(), *ft.Orders[0].ExecutionTime, time.Second)

			return nil
		})

	// Order creation execution
	err := suite.operator.CreateOrder(context.Background(), id, o)
	suite.Require().NoError(err)
}
