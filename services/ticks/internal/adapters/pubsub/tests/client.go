package tests

import (
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

type PubSubClientSuite struct {
	suite.Suite
	Client pubsub.Port
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

func (suite *PubSubClientSuite) TestOnePubTwoSub() {
	as := suite.Require()

	pairSymbol := "symbol2"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ch1, err := suite.Client.Subscribe(pairSymbol)
	as.NoError(err)

	ch2, err := suite.Client.Subscribe(pairSymbol)
	as.NoError(err)

	as.NoError(suite.Client.Publish(t))

	for i := 0; i < 2; i++ {
		select {
		case recvTick := <-ch1:
			as.Equal(t, recvTick)
		case recvTick := <-ch2:
			as.Equal(t, recvTick)
		case <-time.After(1 * time.Second):
			as.FailNow("Timeout")
		}
	}
}

func (suite *PubSubClientSuite) TestCheckClose() {
	as := suite.Require()

	pairSymbol := "symbol3"
	t := tick.Tick{
		Time:       time.Unix(0, 0).UTC(),
		PairSymbol: pairSymbol,
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ch, err := suite.Client.Subscribe(pairSymbol)
	as.NoError(err)

	suite.Client.Close()
	as.Error(suite.Client.Publish(t))

	_, open := <-ch
	suite.False(open)
}
