package ticks

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestListenSuite(t *testing.T) {
	suite.Run(t, new(ListenSuite))
}

type ListenSuite struct {
	suite.Suite
	operator Operator
	vdb      *vdb.MockAdapter
	ps       *pubsub.MockAdapter
	exchange *exchanges.MockAdapter
}

func (suite *ListenSuite) SetupTest() {
	suite.vdb = vdb.NewMockAdapter(gomock.NewController(suite.T()))
	suite.ps = pubsub.NewMockAdapter(gomock.NewController(suite.T()))

	suite.exchange = exchanges.NewMockAdapter(gomock.NewController(suite.T()))
	exchanges := map[string]exchanges.Adapter{"exchange": suite.exchange}

	suite.operator = New(suite.ps, suite.vdb, exchanges)
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
