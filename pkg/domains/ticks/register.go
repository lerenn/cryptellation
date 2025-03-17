package ticks

import (
	"fmt"
	"strings"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal/signals"
	"github.com/lerenn/cryptellation/v1/pkg/temporal/activities"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) RegisterForTicksListeningWorkflow(
	ctx workflow.Context,
	params api.RegisterForTicksListeningWorkflowParams,
) (api.RegisterForTicksListeningWorkflowResults, error) {
	// Check if exchange+pair exists
	if err := wf.checkPairAndExchange(ctx, params.Pair, params.Exchange); err != nil {
		return api.RegisterForTicksListeningWorkflowResults{}, err
	}

	// Send signal-with-start to listen for ticks
	if err := activities.ExecuteSignalWithStart(ctx, activities.SignalWithStartActivityParams{
		SignalName: signals.RegisterToTicksListeningSignalName,
		SignalParams: signals.RegisterToTicksListeningSignalParams{
			CallbackWorkflow: params.Callback,
		},
		WorkflowID:   fmt.Sprintf("%s-%s", strings.ToTitle(params.Exchange), params.Pair),
		WorkflowName: ticksSentryWorkflowName,
		WorkflowParams: ticksSentryWorkflowParams{
			Exchange: params.Exchange,
			Symbol:   params.Pair,
		},
		TaskQueue: api.WorkerTaskQueueName,
	}); err != nil {
		return api.RegisterForTicksListeningWorkflowResults{}, err
	}

	return api.RegisterForTicksListeningWorkflowResults{}, nil
}

func (wf *workflows) checkPairAndExchange(ctx workflow.Context, pair string, exchange string) error {
	// Get exchange info
	result, err := wf.cryptellation.GetExchange(ctx, api.GetExchangeWorkflowParams{
		Name: exchange,
	}, nil)
	if err != nil {
		return err
	}

	// Check if pair exists
	for _, p := range result.Exchange.Pairs {
		if p == pair {
			return nil
		}
	}

	return fmt.Errorf("pair %q doesn't exist for exchange %q", pair, exchange)
}
