package domain

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"

	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type GetStatusSuite struct {
	suite.Suite
	operator     app.ForwardTests
	db           *db.MockPort
	candlesticks *client.MockClient
}

func (suite *GetStatusSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = client.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.candlesticks)
}

func (suite *GetStatusSuite) TestHappyPass() {
	id := uuid.New()

	// Set DB mock expectations
	suite.db.EXPECT().ReadForwardTest(context.Background(), id).
		Return(forwardtest.ForwardTest{
			ID: id,
			Accounts: map[string]account.Account{
				"exchange": {
					Balances: map[string]float64{
						"ETH":  10,
						"BTC":  5,
						"USDT": 1000,
					},
				},
			},
		}, nil)

	// Set Candlesticks mock expectations
	suite.candlesticks.EXPECT().Read(context.Background(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			suite.Require().Equal(payload.Exchange, "exchange")
			suite.Require().Equal(payload.Pair, "ETH-USDT")
			suite.Require().Equal(payload.Period, period.M1)
			suite.Require().NotNil(payload.Start)
			suite.Require().WithinDuration(time.Now().Add(-time.Minute*10), *payload.Start, time.Second)
			suite.Require().NotNil(payload.End)
			suite.Require().WithinDuration(time.Now(), *payload.End, time.Second)
			suite.Require().Equal(payload.Limit, 1)

			return candlestick.NewList("exchange", "ETH-USDT", period.M1).
				MustSet(time.Now(), candlestick.Candlestick{
					Close: 100,
				}), nil
		})
	suite.candlesticks.EXPECT().Read(context.Background(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
			suite.Require().Equal(payload.Exchange, "exchange")
			suite.Require().Equal(payload.Pair, "BTC-USDT")
			suite.Require().Equal(payload.Period, period.M1)
			suite.Require().NotNil(payload.Start)
			suite.Require().WithinDuration(time.Now().Add(-time.Minute*10), *payload.Start, time.Second)
			suite.Require().NotNil(payload.End)
			suite.Require().WithinDuration(time.Now(), *payload.End, time.Second)
			suite.Require().Equal(payload.Limit, 1)

			return candlestick.NewList("exchange", "BTC-USDT", period.M1).
				MustSet(time.Now(), candlestick.Candlestick{
					Close: 1000,
				}), nil
		})

	// Execute getting status
	status, err := suite.operator.GetStatus(context.Background(), id)
	suite.Require().NoError(err)

	suite.Require().Equal(forwardtest.Status{
		Balance: 1000 + 5*1000 + 10*100,
	}, status)
}
