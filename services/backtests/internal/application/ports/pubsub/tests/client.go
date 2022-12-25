package tests

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

func TestPubSubClientSuite(t *testing.T) {
	suite.Run(t, new(PubSubClientSuite))
}

type PubSubClientSuite struct {
	suite.Suite
	Client pubsub.Adapter
}

func (suite *PubSubClientSuite) TestOnePubOneSubObject() {
	backtestID := uint(1)
	ts := time.Unix(60, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	st := status.Status{
		Finished: true,
	}

	ch, err := suite.Client.Subscribe(backtestID)
	suite.Require().NoError(err)

	suite.Require().NoError(suite.Client.Publish(backtestID, event.NewTickEvent(ts, t)))
	select {
	case recvEvent := <-ch:
		suite.checkTick(recvEvent, ts, t)
	case <-time.After(1 * time.Second):
		suite.Require().FailNow("Timeout")
	}

	suite.Require().NoError(suite.Client.Publish(backtestID, event.NewStatusEvent(ts, st)))
	select {
	case recvEvent := <-ch:
		suite.checkEnd(recvEvent, ts, st)
	case <-time.After(1 * time.Second):
		suite.Require().FailNow("Timeout")
	}
}

func (suite *PubSubClientSuite) TestOnePubTwoSub() {
	backtestID := uint(2)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}

	ch1, err := suite.Client.Subscribe(backtestID)
	suite.Require().NoError(err)

	ch2, err := suite.Client.Subscribe(backtestID)
	suite.Require().NoError(err)

	suite.Require().NoError(suite.Client.Publish(backtestID, event.NewTickEvent(ts, t)))

	for i := 0; i < 2; i++ {
		select {
		case recvEvent := <-ch1:
			suite.checkTick(recvEvent, ts, t)
		case recvEvent := <-ch2:
			suite.checkTick(recvEvent, ts, t)
		case <-time.After(1 * time.Second):
			suite.Require().FailNow("Timeout")
		}
	}
}

func (suite *PubSubClientSuite) TestCheckClose() {
	backtestID := uint(3)
	ts := time.Unix(0, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	ch, err := suite.Client.Subscribe(backtestID)
	suite.Require().NoError(err)

	suite.Client.Close()
	suite.Require().Error(suite.Client.Publish(backtestID, event.NewTickEvent(ts, t)))

	_, open := <-ch
	suite.False(open)
}

func (suite *PubSubClientSuite) checkTick(evt event.Event, t time.Time, ti tick.Tick) {
	suite.Require().Equal(event.TypeIsTick, evt.Type)
	suite.Require().Equal(t, evt.Time)
	rt, ok := evt.Content.(tick.Tick)
	suite.Require().True(ok)
	suite.Require().Equal(ti, rt)
}

func (suite *PubSubClientSuite) checkEnd(evt event.Event, t time.Time, st status.Status) {
	suite.Require().Equal(event.TypeIsStatus, evt.Type)
	suite.Require().Equal(t, evt.Time)
	rt, ok := evt.Content.(status.Status)
	suite.Require().True(ok)
	suite.Require().Equal(st, rt)
}
