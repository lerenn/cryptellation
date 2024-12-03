package internal

const TicksSentryWorkflowName = "TicksSentryWorkflow"

type (
	TicksSentryWorkflowParams struct {
		Exchange string
		Symbol   string
	}

	TicksSentryWorkflowResults struct{}
)
