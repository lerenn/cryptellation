package ticks

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/internal/components/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/events"
	"github.com/lerenn/cryptellation/internal/components/ticks/ports/exchanges"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestUnregisterSuite(t *testing.T) {
	suite.Run(t, new(UnregisterSuite))
}

type UnregisterSuite struct {
	suite.Suite
	operator Interface
	vdb      *db.MockPort
	ps       *events.MockPort
	exchange *exchanges.MockPort
}

func (suite *UnregisterSuite) SetupTest() {
	suite.vdb = db.NewMockPort(gomock.NewController(suite.T()))
	suite.ps = events.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))

	suite.operator = New(suite.ps, suite.vdb, suite.exchange)
}

func (suite *UnregisterSuite) setMocksForUnregister() context.Context {
	ctx := context.Background()

	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		DecrementSymbolListenerSubscribers(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(0), nil)

	return ctx
}

func (suite *UnregisterSuite) TestUnregister() {
	ctx := suite.setMocksForUnregister()

	count, err := suite.operator.Unregister(ctx, "exchange", "PAIR_SYMBOL")
	suite.Require().NoError(err)
	suite.Require().Equal(int64(0), count)
}
