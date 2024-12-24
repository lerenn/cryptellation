package ticks

import (
	"errors"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal/signals"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	temporalutils "github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ticksSentryWorkflowName is the name of the TicksSentryWorkflow which is
// a long running workflow that listens to the ticks stream and sends them to
// listeners.
const ticksSentryWorkflowName = "TicksSentryWorkflow"

type (
	// ticksSentryWorkflowParams is the input params for the TicksSentryWorkflow.
	ticksSentryWorkflowParams struct {
		Exchange string
		Symbol   string
	}

	// ticksSentryWorkflowResults is the output results for the TicksSentryWorkflow.
	ticksSentryWorkflowResults struct{}
)

func (wf *workflows) ticksSentryWorkflow(
	ctx workflow.Context,
	params ticksSentryWorkflowParams,
) (ticksSentryWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Listening to ticks",
		"exchange", params.Exchange,
		"symbol", params.Symbol)

	// Get signal channels
	registerSignalChannel := workflow.GetSignalChannel(ctx, signals.RegisterToTicksListeningSignalName)
	unregisterSignalChannel := workflow.GetSignalChannel(ctx, signals.UnregisterFromTicksListeningSignalName)
	newTickReceivedSignalChannel := workflow.GetSignalChannel(ctx, signals.NewTickReceivedSignalName)

	// Start listening to ticks
	cancelListening := wf.sentryStartListeningActivity(ctx, params)

	// Create listeners
	listeners := make(map[string]workflow.Channel)
	handleListenTicksSignals(ctx, listeners, registerSignalChannel, unregisterSignalChannel)

	// Loop over ticks
	var t tick.Tick
	for len(listeners) > 0 {
		// Get next tick
		logger.Debug("Listening to next tick",
			"listeners_count", listeners)
		newTickReceivedSignalChannel.Receive(ctx, &t)

		// Handle new signals
		// TODO: make it async
		handleListenTicksSignals(ctx, listeners, registerSignalChannel, unregisterSignalChannel)

		// Send event to all listeners
		logger.Debug("Sending tick to listeners",
			"tick", t,
			"listeners_count", len(listeners))
		keys := workflow.DeterministicKeys(listeners)
		for _, k := range keys {
			_ = listeners[k].SendAsync(t)
		}
	}

	// Cancel listening and cleanup signals
	logger.Debug("No more listeners, cancel listening")
	cancelListening()

	// Cleanup remaining signals
	// TODO: clean up new ticks signals
	// TODO: clean up unregister signals
	// TODO: clean up register signals and trigger a new workflow if needed

	logger.Info("Stop listening to ticks",
		"exchange", params.Exchange,
		"symbol", params.Symbol)

	return ticksSentryWorkflowResults{}, nil
}

func (wf *workflows) sentryStartListeningActivity(
	ctx workflow.Context,
	params ticksSentryWorkflowParams,
) func() {
	// Set activity options
	// TODO: Improve this
	activityOptions := exchanges.DefaultActivityOptions()
	activityOptions.ScheduleToCloseTimeout = 365 * 24 * time.Hour
	activityOptions.StartToCloseTimeout = 365 * 24 * time.Hour
	activityOptions.HeartbeatTimeout = time.Second
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute activity with cancel
	ctx, cancelActivity := workflow.WithCancel(ctx)
	_ = workflow.ExecuteActivity(
		ctx, wf.exchanges.ListenSymbolActivity, exchanges.ListenSymbolParams{
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
	logger.Debug("Handling signals",
		"listeners_count", len(listeners))

	// Handle register signals
	var registerParams signals.RegisterToTicksListeningSignalParams
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
	var unregisterParams signals.UnregisterFromTicksListeningSignalParams
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
	callback temporalutils.CallbackWorkflow,
	listeners map[string]workflow.Channel,
) func(ctx workflow.Context) {
	ch := listeners[callback.Name]

	// Options
	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                callback.TaskQueueName,            // Execute in the client queue
		ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON, // Do not close if the parent workflow closes
		WorkflowExecutionTimeout: time.Second * 30,                  // Timeout if the child workflow does not complete
	}

	// Check if the timeout is set
	if callback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = callback.ExecutionTimeout
	}

	// Return function that will send signal to new child workflow
	return func(ctx workflow.Context) {
		var t tick.Tick
		logger := workflow.GetLogger(ctx)

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
			if err != nil {
				if errors.As(err, &timeoutErr) {
					logger.Debug("Listener has timed out, exiting",
						"callback", callback.Name)
					break
				}

				logger.Error("Listener has errored, continuing",
					"error", err,
					"callback", callback.Name)
			}
		}

		// Remove listener as it has been in error or stopped.
		logger.Debug("Removing listener",
			"callback", callback.Name)
		delete(listeners, callback.Name)
	}
}
