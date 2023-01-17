package backtests

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	candlesticks "github.com/digital-feather/cryptellation/services/candlesticks/clients/go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestCreateOrderSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSuite))
}

type CreateOrderSuite struct {
	suite.Suite
	operator     Operator
	db           *db.MockAdapter
	pubsub       *pubsub.MockAdapter
	candlesticks *candlesticks.MockInterfacer
}

func (suite *CreateOrderSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.pubsub = pubsub.NewMockAdapter(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockInterfacer(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.pubsub, suite.candlesticks)
}

func (suite *CreateOrderSuite) TestHappyPass() {
	// TODO
}
