package tests

import (
	"time"

	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/pubsub"
	"github.com/digital-feather/cryptellation/pkg/tick"
	"github.com/stretchr/testify/suite"
)

type PubSubClientSuite struct {
	suite.Suite
	Client pubsub.Port
}

func (suite *PubSubClientSuite) TearDownTest() {
	suite.Client.Close()
}

func (suite *PubSubClientSuite) TestOnePubOneSubObject() {
	as := suite.Require()

	pairSymbol := "symbol1"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	ch, err := suite.Client.Subscribe(pairSymbol)
	as.NoError(err)

	as.NoError(suite.Client.Publish(t))
	select {
	case recvTick := <-ch:
		as.Equal(t, recvTick)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}
}
