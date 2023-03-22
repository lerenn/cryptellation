package app

import (
	"testing"

	"github.com/digital-feather/cryptellation/clients/go/mock"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestGetAccountsSuite(t *testing.T) {
	suite.Run(t, new(GetAccountsSuite))
}

type GetAccountsSuite struct {
	suite.Suite
	operator     Controller
	db           *db.MockAdapter
	Events       *events.MockAdapter
	candlesticks *mock.MockCandlesticks
}

func (suite *GetAccountsSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.Events = events.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *GetAccountsSuite) TestHappyPass() {
	// TODO
}
