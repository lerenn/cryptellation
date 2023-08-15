package ticks

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/internal/core/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/core/ticks/ports/events"
	"github.com/lerenn/cryptellation/internal/core/ticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

type RegisterSuite struct {
	suite.Suite
	operator Interface
	vdb      *db.MockPort
	ps       *events.MockPort
	exchange *exchanges.MockPort
}

func (suite *RegisterSuite) SetupTest() {
	suite.vdb = db.NewMockPort(gomock.NewController(suite.T()))
	suite.ps = events.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))

	suite.operator = New(suite.ps, suite.vdb, suite.exchange)
}

func (suite *RegisterSuite) setMocksForFirstRegister(ctx context.Context) (chan tick.Tick, func(), *sync.WaitGroup) {
	ch := make(chan tick.Tick, 10)

	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		IncrementSymbolListenerSubscribers(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(1), nil)

	// Set call to exchange to listen to symbol
	suite.exchange.EXPECT().
		ListenSymbol("exchange", "PAIR_SYMBOL").
		Return(ch, make(chan struct{}, 10), nil)

	// Set call to Events when receving a tick for the exchange
	wg := sync.WaitGroup{}
	suite.ps.EXPECT().Publish(tick.Tick{
		Time:       time.Unix(60, 0),
		PairSymbol: "SYMBOL",
		Price:      2.0,
		Exchange:   "EXCHANGE",
	}).DoAndReturn(func(tick tick.Tick) error {
		wg.Done()
		return nil
	})
	wg.Add(1)

	// Set call to Events when closing the goroutine automatically
	closeWaitGroup := sync.WaitGroup{}
	suite.ps.EXPECT().Close().Do(func() {
		closeWaitGroup.Done()
	})
	closeWaitGroup.Add(1)

	return ch, func() {
		close(ch)
		closeWaitGroup.Wait()
	}, &wg
}

func (suite *RegisterSuite) TestFirstRegister() {
	ctx := context.Background()
	fromExchangeChan, cleanup, wg := suite.setMocksForFirstRegister(ctx)
	defer cleanup()

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
	fromExchangeChan <- t

	// Wait for tick to be arrived
	wg.Wait()
}

func (suite *RegisterSuite) setMocksForSecondRegister() context.Context {
	ctx := context.Background()

	// Set call to database for checking existing listener, and return the new count
	suite.vdb.EXPECT().
		IncrementSymbolListenerSubscribers(ctx, "exchange", "PAIR_SYMBOL").
		Return(int64(2), nil)

	// Nothing more should happen

	return ctx
}

func (suite *RegisterSuite) TestSecondRegister() {
	ctx := suite.setMocksForSecondRegister()

	// Register to the application
	count, err := suite.operator.Register(ctx, "exchange", "PAIR_SYMBOL")

	// Check return
	suite.Require().NoError(err)
	suite.Require().Equal(int64(2), count)
}

// TODO: find a way to check real closure when no listener
