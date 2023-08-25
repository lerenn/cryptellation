package events

import (
	context "context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

type EventsClientSuite struct {
	suite.Suite
	Client Port
}

func (suite *EventsClientSuite) TearDownTest() {
	suite.Client.Close(context.Background())
}

func (suite *EventsClientSuite) TestOnePubOneSubObject() {
	as := suite.Require()

	pairSymbol := "symbol1"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	ch, err := suite.Client.Subscribe(context.Background(), pairSymbol)
	as.NoError(err)

	as.NoError(suite.Client.Publish(context.Background(), t))
	select {
	case recvTick := <-ch:
		as.Equal(t, recvTick)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}
}
