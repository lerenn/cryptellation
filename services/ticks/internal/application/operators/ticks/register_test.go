package ticks

import (
	"context"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

type RegisterSuite struct {
	suite.Suite
	operator Operator
	vdb      *vdb.MockAdapter
	ps       *pubsub.MockAdapter
	exchange *exchanges.MockAdapter
}

func (suite *RegisterSuite) SetupTest() {
	suite.vdb = vdb.NewMockAdapter(gomock.NewController(suite.T()))
	suite.ps = pubsub.NewMockAdapter(gomock.NewController(suite.T()))

	suite.exchange = exchanges.NewMockAdapter(gomock.NewController(suite.T()))
	exchanges := map[string]exchanges.Adapter{"exchange": suite.exchange}

	suite.operator = New(suite.ps, suite.vdb, exchanges)
}

func (suite *RegisterSuite) setMocksForFirstRegister(ctx context.Context, ch chan tick.Tick, stopCh chan struct{}) {
	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		IncrementSymbolListenerCount(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(1), nil)

	// Set call to exchange to listen to symbol
	suite.exchange.EXPECT().
		ListenSymbol("PAIR_SYMBOL").
		Return(ch, stopCh, nil)

	// Set call to pubsub when receving a tick for the exchange
	suite.ps.EXPECT().Publish(tick.Tick{
		Time:       time.Unix(60, 0),
		PairSymbol: "SYMBOL",
		Price:      2.0,
		Exchange:   "EXCHANGE",
	}).Return(nil)

	// Set call to pubsub when closing the goroutine automatically
	suite.ps.EXPECT().Close()
}

func (suite *RegisterSuite) TestFirstRegister() {
	ctx := context.Background()
	ch := make(chan tick.Tick, 1)
	defer close(ch)
	stopCh := make(chan struct{}, 1)
	defer close(stopCh)
	suite.setMocksForFirstRegister(ctx, ch, stopCh)

	// Register to the application
	count, err := suite.operator.Register(ctx, "exchange", "PAIR_SYMBOL")

	// Check return
	suite.Require().NoError(err)
	suite.Require().Equal(int64(1), count)

	// Simulate sending a tick from the exchange
	t := tick.Tick{
		Time:       time.Unix(60, 0),
		PairSymbol: "SYMBOL",
		Price:      2.0,
		Exchange:   "EXCHANGE",
	}
	ch <- t
	time.Sleep(time.Millisecond) // Wait for channel to propagate
}

func (suite *RegisterSuite) setMocksForSecondRegister(ctx context.Context) {
	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		IncrementSymbolListenerCount(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(2), nil)

	// Nothing more should happen
}

func (suite *RegisterSuite) TestSecondRegister() {
	ctx := context.Background()
	suite.setMocksForSecondRegister(ctx)

	// Register to the application
	count, err := suite.operator.Register(ctx, "exchange", "PAIR_SYMBOL")

	// Check return
	suite.Require().NoError(err)
	suite.Require().Equal(int64(2), count)
}

// TODO: find a way to check real closure when no listener
