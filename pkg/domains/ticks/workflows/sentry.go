package workflows

import (
	"errors"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) TicksSentryWorkflow(
	ctx workflow.Context,
	params internal.TicksSentryWorkflowParams,
) (internal.TicksSentryWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Listening to ticks",
		"exchange", params.Exchange,
		"symbol", params.Symbol)

	// Get signal channels
	registerSignalChannel := workflow.GetSignalChannel(ctx, internal.RegisterToTicksListeningSignalName)
	unregisterSignalChannel := workflow.GetSignalChannel(ctx, internal.UnregisterFromTicksListeningSignalName)
	newTickReceivedSignalChannel := workflow.GetSignalChannel(ctx, internal.NewTickReceivedSignalName)

	// Start listening to ticks
	cancelListening := wf.sentryStartListeningActivity(ctx, params)

	// Create listeners
	listeners := make(map[string]workflow.Channel)
	handleListenTicksSignals(ctx, listeners, registerSignalChannel, unregisterSignalChannel)

	// Loop over ticks
	var t tick.Tick
	for len(listeners) > 0 {
		// Get next tick
		newTickReceivedSignalChannel.Receive(ctx, &t)

		// Handle new signals
		// TODO: make it async
		handleListenTicksSignals(ctx, listeners, registerSignalChannel, unregisterSignalChannel)

		// Send event to all listeners
		keys := workflow.DeterministicKeys(listeners)
		for _, k := range keys {
			_ = listeners[k].SendAsync(t)
		}
	}

	// Cancel listening and cleanup signals
	cancelListening()

	// Cleanup remaining signals
	// TODO: clean up new ticks signals
	// TODO: clean up unregister signals
	// TODO: clean up register signals and trigger a new workflow if needed

	logger.Info("Stop listening to ticks",
		"exchange", params.Exchange,
		"symbol", params.Symbol)

	return internal.TicksSentryWorkflowResults{}, nil
}

func (wf *workflows) sentryStartListeningActivity(
	ctx workflow.Context,
	params internal.TicksSentryWorkflowParams,
) func() {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 365 * 24 * time.Hour,
		HeartbeatTimeout:    time.Second,
	}
	aCtx := workflow.WithActivityOptions(ctx, activityOptions)
	aCtx, cancelActivity := workflow.WithCancel(aCtx)
	_ = workflow.ExecuteActivity(
		aCtx, wf.exchanges.ListenSymbolActivity, exchanges.ListenSymbolParams{
			ParentWorkflowID: workflow.GetInfo(ctx).WorkflowExecution.ID,
			Exchange:         params.Exchange,
			Symbol:           params.Symbol,
		})
	return cancelActivity
}

func handleListenTicksSignals(
	ctx workflow.Context,
	listeners map[string]workflow.Channel,
	registerSignalChannel, unregisterSignalChannel workflow.ReceiveChannel,
) {
	logger := workflow.GetLogger(ctx)

	// Handle register signals
	var registerParams internal.RegisterToTicksListeningSignalParams
	for detected := true; detected; {
		detected = registerSignalChannel.ReceiveAsync(&registerParams)
		if detected {
			logger.Info("Received register signal",
				"params", registerParams)
			listeners[registerParams.CallbackWorkflow.Name] = workflow.NewBufferedChannel(ctx, 0)
			workflow.Go(ctx, sendToTickListenerRoutine(registerParams.CallbackWorkflow, listeners))
		}
	}

	// Handle unregister signals
	var unregisterParams internal.UnregisterFromTicksListeningSignalParams
	for detected := true; detected; {
		detected = unregisterSignalChannel.ReceiveAsync(&unregisterParams)
		if detected {
			logger.Info("Received unregister signal",
				"params", unregisterParams)
			delete(listeners, unregisterParams.CallbackWorkflowName)
		}
	}
}

func sendToTickListenerRoutine(
	callback api.ListenToTicksCallbackWorkflow,
	listeners map[string]workflow.Channel,
) func(ctx workflow.Context) {
	ch := listeners[callback.Name]

	// Options
	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                callback.TaskQueueName,            // Execute in the client queue
		ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON, // Do not close if the parent workflow closes
		WorkflowExecutionTimeout: time.Second * 30,                  // Timeout after 30 seconds if the child workflow does not complete
	}

	// Check if the timeout is set
	if callback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = callback.ExecutionTimeout
	}

	// Return function that will send signal to new child workflow
	return func(ctx workflow.Context) {
		var t tick.Tick

		ctx = workflow.WithChildOptions(ctx, opts)
		for {
			// Receive next event
			ch.Receive(ctx, &t)

			// Start a new child workflow
			err := workflow.ExecuteChildWorkflow(ctx, callback.Name, api.ListenToTicksCallbackWorkflowParams{
				Tick: t,
			}).Get(ctx, nil)

			// Only stop if time-out
			var timeoutErr *temporal.TimeoutError
			if err != nil && errors.As(err, &timeoutErr) {
				break
			}
		}

		// Remove listener as it has been in error or stopped.
		delete(listeners, callback.Name)
	}
}
