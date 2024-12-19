package direct

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	temporalutils "github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func (c client) ListenToTicks(
	ctx context.Context,
	exchange, pair string,
	callback func(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error,
) error {
	// Create variables
	tq := fmt.Sprintf("ListenTicks-%s", uuid.New().String())
	workflowName := tq

	// Create temporary worker
	w := worker.New(c.Temporal(), tq, worker.Options{})
	w.RegisterWorkflowWithOptions(callback, workflow.RegisterOptions{
		Name: workflowName,
	})

	// Start worker
	go func() {
		if err := w.Run(nil); err != nil {
			panic(err) // TODO: Handle error by returning it if there is an error
		}
	}()
	defer w.Stop()

	// Listen to ticks
	_, err := c.Client.ListenToTicks(ctx,
		api.RegisterForTicksListeningWorkflowParams{
			Exchange: exchange,
			Pair:     pair,
			Callback: temporalutils.CallbackWorkflow{
				Name:          workflowName,
				TaskQueueName: tq,
			},
		})
	if err != nil {
		return err
	}

	// Wait for interrupt
	<-ctx.Done()

	return nil
}
