package events

import (
	context "context"
	"time"

	event "github.com/lerenn/cryptellation/pkg/models/event"

	tick "github.com/lerenn/cryptellation/ticks/pkg/tick"

	"github.com/stretchr/testify/suite"
)

type EventsClientSuite struct {
	suite.Suite
	Client Port
}

func (suite *EventsClientSuite) TearDownTest() {
	suite.Client.Close(context.Background())
}

func (suite *EventsClientSuite) TestPublishTick() {
	as := suite.Require()

	t := tick.Tick{
		Exchange: "exchange1",
		Pair:     "symbol1",
		Time:     time.Unix(0, 0).UTC(),
		Price:    float64(time.Now().UnixNano()),
	}

	ch, err := suite.Client.SubscribeToTicks(context.Background(), event.TickSubscription{
		Exchange: t.Exchange,
		Pair:     t.Pair,
	})
	as.NoError(err)

	as.NoError(suite.Client.PublishTick(context.Background(), t))
	select {
	case recvTick := <-ch:
		as.Equal(t, recvTick)
	case <-time.After(1 * time.Second):
		as.FailNow("Timeout")
	}
}
