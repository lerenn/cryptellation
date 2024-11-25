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
	// ServiceInfoParams contains the parameters of the service info workflow.
	ServiceInfoParams struct{}

	// ServiceInfoResult contains the result of the service info workflow.
	ServiceInfoResult struct {
		Version string
	}
)
