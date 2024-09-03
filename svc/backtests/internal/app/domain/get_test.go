package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type GetSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *GetSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *GetSuite) TestHappyPass() {
	// Set DB mock expectations
	bt := backtest.Backtest{
		ID: uuid.New(),
	}
	suite.db.EXPECT().
		ReadBacktest(gomock.Any(), bt.ID).
		Return(bt, nil)

	// Get backtest
	b, err := suite.operator.Get(context.Background(), bt.ID)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(bt, b)
}
