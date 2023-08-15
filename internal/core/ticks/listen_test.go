package ticks

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/internal/core/ticks/ports/db"
	"github.com/lerenn/cryptellation/internal/core/ticks/ports/events"
	"github.com/lerenn/cryptellation/internal/core/ticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestListenSuite(t *testing.T) {
	suite.Run(t, new(ListenSuite))
}

type ListenSuite struct {
	suite.Suite
	operator Interface
	vdb      *db.MockPort
	ps       *events.MockPort
	exchange *exchanges.MockPort
}

func (suite *ListenSuite) SetupTest() {
	suite.vdb = db.NewMockPort(gomock.NewController(suite.T()))
	suite.ps = events.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))

	suite.operator = New(suite.ps, suite.vdb, suite.exchange)
}

func (suite *ListenSuite) setMocksForHappyPath() chan tick.Tick {
	ch := make(chan tick.Tick, 1)

	// Set the expected call for subscribing to the messages
	suite.ps.EXPECT().Subscribe("SYMBOL").Return(ch, nil)

	return ch
}

func (suite *ListenSuite) TestHappyPass() {
	ch := suite.setMocksForHappyPath()

	// Make the call
	rch, err := suite.operator.Listen("EXCHANGE", "SYMBOL")

	// Check returned values
	suite.Require().NoError(err)

	// Send a tick from mock perspective
	t := tick.Tick{
		Time:       time.Unix(60, 0),
		PairSymbol: "SYMBOL",
		Price:      2.0,
		Exchange:   "EXCHANGE",
	}
	ch <- t

	// Check reception
	rt, ok := <-rch
	suite.Require().True(ok)
	suite.Require().Equal(t, rt)
}
