package events

import (
	context "context"

	"github.com/stretchr/testify/suite"
)

type EventsClientSuite struct {
	suite.Suite
	Client Port
}

func (suite *EventsClientSuite) TearDownTest() {
	suite.Client.Close(context.Background())
}

// Disable test until ticks rework
// func (suite *EventsClientSuite) TestOnePubOneSubObject() {
// 	as := suite.Require()

// 	pair := "symbol1"
// 	t := tick.Tick{
// 		Time:     time.Unix(0, 0).UTC(),
// 		Pair:     pair,
// 		Price:    float64(time.Now().UnixNano()),
// 		Exchange: "exchange",
// 	}
// 	ch, err := suite.Client.Subscribe(context.Background(), pair)
// 	as.NoError(err)

// 	as.NoError(suite.Client.Publish(context.Background(), t))
// 	select {
// 	case recvTick := <-ch:
// 		as.Equal(t, recvTick)
// 	case <-time.After(1 * time.Second):
// 		as.FailNow("Timeout")
// 	}
// }
