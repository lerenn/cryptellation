package app

import (
	"testing"

	"github.com/digital-feather/cryptellation/clients/go/mock"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestGetOrdersSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersSuite))
}

type GetOrdersSuite struct {
	suite.Suite
	operator     Controller
	db           *db.MockAdapter
	Events       *events.MockAdapter
	candlesticks *mock.MockCandlesticks
}

func (suite *GetOrdersSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.Events = events.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *GetOrdersSuite) TestHappyPass() {
	// TODO
}
