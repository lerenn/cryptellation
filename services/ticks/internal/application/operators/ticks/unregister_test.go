package ticks

import (
	"context"
	"testing"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestUnregisterSuite(t *testing.T) {
	suite.Run(t, new(UnregisterSuite))
}

type UnregisterSuite struct {
	suite.Suite
	operator Operator
	vdb      *vdb.MockAdapter
	ps       *pubsub.MockAdapter
	exchange *exchanges.MockAdapter
}

func (suite *UnregisterSuite) SetupTest() {
	suite.vdb = vdb.NewMockAdapter(gomock.NewController(suite.T()))
	suite.ps = pubsub.NewMockAdapter(gomock.NewController(suite.T()))

	suite.exchange = exchanges.NewMockAdapter(gomock.NewController(suite.T()))
	exchanges := map[string]exchanges.Adapter{"exchange": suite.exchange}

	suite.operator = New(suite.ps, suite.vdb, exchanges)
}

func (suite *UnregisterSuite) setMocksForUnregister(ctx context.Context) {
	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		DecrementSymbolListenerCount(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(0), nil)
}

func (suite *UnregisterSuite) TestUnregister() {
	ctx := context.Background()
	suite.setMocksForUnregister(ctx)
}
