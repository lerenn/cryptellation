package test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func (suite *EndToEndSuite) TestListenToTicks() {
	exchange := "binance"
	pair := "BTC-USDT"
	count := 0

	// GIVEN a Temporal worker

	// Create the worker
	tq := fmt.Sprintf("E2E-Run-%s", uuid.New().String())
	w := worker.New(suite.client.Temporal(), tq, worker.Options{})
	w.RegisterWorkflowWithOptions(func(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
		suite.Require().Equal(exchange, params.Tick.Exchange)
		suite.Require().Equal(pair, params.Tick.Pair)
		count++
		return nil
	}, workflow.RegisterOptions{
		Name: tq,
	})

	// Start worker
	irq := worker.InterruptCh()
	go w.Run(irq)

	// WHEN registering for ticks listening

	_, err := suite.client.ListenToTicks(context.Background(),
		api.RegisterForTicksListeningWorkflowParams{
			Exchange: exchange,
			Pair:     pair,
			CallbackWorkflow: api.ListenToTicksCallbackWorkflow{
				Name:          tq,
				TaskQueueName: tq,
			},
		})

	// THEN no error is returned

	suite.Require().NoError(err)

	// AND the count is increased after a while

	suite.Eventually(func() bool {
		return count > 0
	}, time.Minute, time.Second,
		"count should be greater than 0")

	// Stop worker
	w.Stop()
}
