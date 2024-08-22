package domain

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

type ListSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *ListSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *ListSuite) TestHappyPath() {
	bt1 := backtest.Backtest{
		ID: uuid.New(),
	}
	bt2 := backtest.Backtest{
		ID: uuid.New(),
	}

	// Set DB mock expectations
	suite.db.EXPECT().
		ListBacktests(context.Background()).
		Return([]backtest.Backtest{bt1, bt2}, nil)

	// List backtests
	bts, err := suite.operator.List(context.Background())
	suite.Require().NoError(err)

	// Check response
	suite.Require().Len(bts, 2)
	suite.Require().Equal(bt1, bts[0])
	suite.Require().Equal(bt2, bts[1])
}
