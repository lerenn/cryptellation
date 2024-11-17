package api

const (
	// WorkerTaskQueueName is the name of the task queue for the cryptellation worker.
	WorkerTaskQueueName = "CryptellationTaskQueue"
)

const (
	// ServiceInfoWorkflowName is the name of the workflow to get the service info.
	ServiceInfoWorkflowName = "ServiceInfoWorkflow"
)

type (
	// ServiceInfoWorkflowResult contains the result of the service info workflow.
	ServiceInfoWorkflowResult struct {
		Version string
	}
)
