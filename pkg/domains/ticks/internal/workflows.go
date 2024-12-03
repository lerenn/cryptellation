package internal

// TicksSentryWorkflowName is the name of the TicksSentryWorkflow which is
// a long running workflow that listens to the ticks stream and sends them to
// listeners.
const TicksSentryWorkflowName = "TicksSentryWorkflow"

type (
	// TicksSentryWorkflowParams is the input params for the TicksSentryWorkflow.
	TicksSentryWorkflowParams struct {
		Exchange string
		Symbol   string
	}

	// TicksSentryWorkflowResults is the output results for the TicksSentryWorkflow.
	TicksSentryWorkflowResults struct{}
)
