package bot

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func RegisterWorkflows(w worker.Worker, taskQueue string, id uuid.UUID, bot Bot) api.Callbacks {
	// Register OnInitCallback
	onInitCallbackWorkflowName := fmt.Sprintf("OnInit-%s", id.String())
	w.RegisterWorkflowWithOptions(bot.OnInit, workflow.RegisterOptions{
		Name: onInitCallbackWorkflowName,
	})

	// Register OnNewPricesCallback
	onNewPricesCallbackWorkflowName := fmt.Sprintf("OnNewPrices-%s", id.String())
	w.RegisterWorkflowWithOptions(bot.OnNewPrices, workflow.RegisterOptions{
		Name: onNewPricesCallbackWorkflowName,
	})

	// Register OnExitCallback
	onExitCallbackWorkflowName := fmt.Sprintf("OnExit-%s", id.String())
	w.RegisterWorkflowWithOptions(bot.OnExit, workflow.RegisterOptions{
		Name: onExitCallbackWorkflowName,
	})

	return api.Callbacks{
		OnInitCallback: temporal.CallbackWorkflow{
			Name:          onInitCallbackWorkflowName,
			TaskQueueName: taskQueue,
		},
		OnNewPricesCallback: temporal.CallbackWorkflow{
			Name:          onNewPricesCallbackWorkflowName,
			TaskQueueName: taskQueue,
		},
		OnExitCallback: temporal.CallbackWorkflow{
			Name:          onExitCallbackWorkflowName,
			TaskQueueName: taskQueue,
		},
	}
}
