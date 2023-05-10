package backtests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lerenn/cryptellation/clients/go/mock"
	"github.com/lerenn/cryptellation/services/backtests/io/db"
	"github.com/lerenn/cryptellation/services/backtests/io/events"
	"github.com/stretchr/testify/suite"
)

func TestCreateOrderSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSuite))
}

type CreateOrderSuite struct {
	suite.Suite
	operator     Interface
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *mock.MockCandlesticks
}

func (suite *CreateOrderSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *CreateOrderSuite) TestHappyPass() {
	// TODO
}
