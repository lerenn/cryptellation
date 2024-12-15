package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

func (suite *EndToEndSuite) TestListenToTicks() {
	exchange := "binance"
	pair := "BTC-USDT"
	count := 0

	// WHEN registering for ticks listening

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var listenErr error
	go func() {
		listenErr = suite.client.ListenToTicks(ctx, exchange, pair,
			func(_ workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
				suite.Require().Equal(exchange, params.Tick.Exchange)
				suite.Require().Equal(pair, params.Tick.Pair)
				count++
				return nil
			})
	}()

	// THEN the count is increased after a while

	suite.Eventually(func() bool {
		return count > 0
	}, 10*time.Minute, time.Second,
		"count should be greater than 0")

	// WHEN cancelling the context

	cancel()

	// THEN no error is returned

	suite.Require().NoError(listenErr)
}
