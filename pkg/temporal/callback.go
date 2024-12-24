package temporal

import "time"

// CallbackWorkflow is the parameters of a callback workflow.
type CallbackWorkflow struct {
	Name             string
	TaskQueueName    string
	ExecutionTimeout time.Duration
}
