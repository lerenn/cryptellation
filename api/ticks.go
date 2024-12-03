package api

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	"go.temporal.io/sdk/worker"
)

const (
	// RegisterForTicksListeningWorkflowName is the name of the workflow to register
	// for ticks reception through a callback workflow.
	RegisterForTicksListeningWorkflowName = "RegisterForTicksListeningWorkflow"
)

type (
	ListenToTicksCallbackWorkflow struct {
		Name             string
		TaskQueueName    string
		ExecutionTimeout time.Duration
	}

	// RegisterForTicksListeningWorkflowParams is the parameters of the
	// RegisterForTicksListening workflow.
	RegisterForTicksListeningWorkflowParams struct {
		Exchange         string
		Pair             string
		CallbackWorkflow ListenToTicksCallbackWorkflow
	}

	// ListenToTicksCallbackWorkflowParams is the parameters of the
	// RegisterForTicksListening callback workflow.
	ListenToTicksCallbackWorkflowParams struct {
		Tick tick.Tick
	}

	// RegisterForTicksListeningWorkflowResults is the results of the
	// RegisterForTicksListening workflow.
	RegisterForTicksListeningWorkflowResults struct {
		Worker worker.Worker
	}
)

const (
	// UnregisterFromTicksListeningWorkflowName is the name of the workflow to register
	// for ticks reception through a callback workflow.
	UnregisterFromTicksListeningWorkflowName = "UnregisterFromTicksListeningWorkflow"
)

type (
	// UnregisterFromTicksListeningWorkflowParams is the parameters of the
	// UnregisterFromTicksListening workflow.
	UnregisterFromTicksListeningWorkflowParams struct {
		CallbackWorkflowName string
	}

	// UnregisterFromTicksListeningWorkflowResults is the results of the
	// UnregisterFromTicksListening workflow.
	UnregisterFromTicksListeningWorkflowResults struct{}
)
