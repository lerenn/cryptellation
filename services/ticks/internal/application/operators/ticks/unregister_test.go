package ticks

import (
	"context"
	"testing"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestUnregisterSuite(t *testing.T) {
	suite.Run(t, new(UnregisterSuite))
}

type UnregisterSuite struct {
	suite.Suite
	operator Operator
	vdb      *db.MockAdapter
	ps       *pubsub.MockAdapter
	exchange *exchanges.MockAdapter
}

func (suite *UnregisterSuite) SetupTest() {
	suite.vdb = db.NewMockAdapter(gomock.NewController(suite.T()))
	suite.ps = pubsub.NewMockAdapter(gomock.NewController(suite.T()))

	suite.exchange = exchanges.NewMockAdapter(gomock.NewController(suite.T()))
	exchanges := map[string]exchanges.Adapter{"exchange": suite.exchange}

	suite.operator = New(suite.ps, suite.vdb, exchanges)
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
