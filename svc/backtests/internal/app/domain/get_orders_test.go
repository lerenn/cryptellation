package domain

import (
	"testing"

	"cryptellation/svc/backtests/internal/app"
	"cryptellation/svc/backtests/internal/app/ports/db"
	"cryptellation/svc/backtests/internal/app/ports/events"

	candlesticks "cryptellation/svc/candlesticks/clients/go"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetOrdersSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersSuite))
}

type GetOrdersSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *GetOrdersSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *GetOrdersSuite) TestHappyPass() {
	// TODO
}
