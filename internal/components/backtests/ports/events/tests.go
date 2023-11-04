package events

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

func TestEventsClientSuite(t *testing.T) {
	suite.Run(t, new(EventsClientSuite))
}

type EventsClientSuite struct {
	suite.Suite
	Client Port
}

func (suite *EventsClientSuite) TestOnePubOneSubObject() {
	backtestID := uint(1)
	ts := time.Unix(60, 0).UTC()
	t := tick.Tick{
		PairSymbol: "BTC-USDC",
		Price:      float64(time.Now().UnixNano()),
		Exchange:   "exchange",
	}
	st := event.Status{
		Finished: true,
	}

	ch, err := suite.Client.Subscribe(context.Background(), backtestID)
	suite.Require().NoError(err)

	suite.Require().NoError(suite.Client.Publish(context.Background(), backtestID, event.NewTickEvent(ts, t)))
	select {
	case recvEvent := <-ch:
		suite.checkTick(recvEvent, ts, t)
	case <-time.After(1 * time.Second):
		suite.Require().FailNow("Timeout")
	}

	suite.Require().NoError(suite.Client.Publish(context.Background(), backtestID, event.NewStatusEvent(ts, st)))
	select {
	case recvEvent := <-ch:
		suite.checkEnd(recvEvent, ts, st)
	case <-time.After(1 * time.Second):
		suite.Require().FailNow("Timeout")
	}
}

func (suite *EventsClientSuite) checkTick(evt event.Event, t time.Time, ti tick.Tick) {
	suite.Require().Equal(event.TypeIsTick, evt.Type)
	suite.Require().Equal(t, evt.Time)
	rt, ok := evt.Content.(tick.Tick)
	suite.Require().True(ok)
	suite.Require().Equal(ti, rt)
}

func (suite *EventsClientSuite) checkEnd(evt event.Event, t time.Time, st event.Status) {
	suite.Require().Equal(event.TypeIsStatus, evt.Type)
	suite.Require().Equal(t, evt.Time)
	rt, ok := evt.Content.(event.Status)
	suite.Require().True(ok)
	suite.Require().Equal(st, rt)
}
