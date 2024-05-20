package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/event"
)

func (suite *EndToEndSuite) TestListenTicks() {
	// WHEN subscribing to ticks
	cancelableCtx, cancel := context.WithCancel(context.Background())
	ch, err := suite.client.SubscribeToTicks(cancelableCtx, event.TickSubscription{
		Exchange: "binance",
		Pair:     "ETH-USDT",
	})
	defer cancel()

	// THEN no error occurs
	suite.Require().NoError(err)

	// AND ticks are received
	select {
	case tick := <-ch:
		suite.Require().NotEmpty(tick)
	case <-time.After(10 * time.Second):
		suite.Fail("No tick received")
	}
}
