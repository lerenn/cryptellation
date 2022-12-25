package backtests

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestGetOrdersSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersSuite))
}

type GetOrdersSuite struct {
	suite.Suite
	operator     Operator
	db           *db.MockAdapter
	pubsub       *pubsub.MockAdapter
	candlesticks *MockClient
}

func (suite *GetOrdersSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
}

func (suite *GetOrdersSuite) TestHappyPass() {
	// TODO
}
