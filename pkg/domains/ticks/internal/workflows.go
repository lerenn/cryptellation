package internal

const ListenToTicksWorkflowName = "ListenToTicksWorkflow"

type (
	ListenToTicksWorkflowParams struct {
		Exchange string
		Symbol   string
	}

	ListenToTicksWorkflowResults struct{}
)
