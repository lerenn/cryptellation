package api

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
)

const (
	// RegisterForTicksListeningWorkflowName is the name of the workflow to register
	// for ticks reception through a callback workflow.
	RegisterForTicksListeningWorkflowName = "RegisterForTicksListeningWorkflow"
)

type (
	RegisterForTicksCallbackWorkflow struct {
		Name             string
		TaskQueueName    string
		ExecutionTimeout time.Duration
	}

	// RegisterForTicksListeningWorkflowParams is the parameters of the
	// RegisterForTicksReception workflow.
	RegisterForTicksListeningWorkflowParams struct {
		Exchange         string
		Pair             string
		CallbackWorkflow RegisterForTicksCallbackWorkflow
	}

	// RegisterForTicksListeningWorkflowCallbackParams is the parameters of the
	// RegisterForTicksReception callback workflow.
	RegisterForTicksListeningWorkflowCallbackParams struct {
		Tick tick.Tick
	}

	// RegisterForTicksListeningWorkflowResults is the results of the
	// RegisterForTicksReception workflow.
	RegisterForTicksListeningWorkflowResults struct{}
)

const (
	// UnregisterFromTicksListeningWorkflowName is the name of the workflow to register
	// for ticks reception through a callback workflow.
	UnregisterFromTicksListeningWorkflowName = "UnregisterFromTicksListeningWorkflow"
)

type (
	// UnregisterFromTicksListeningWorkflowParams is the parameters of the
	// UnregisterFromTicksReception workflow.
	UnregisterFromTicksListeningWorkflowParams struct {
		CallbackWorkflowName string
	}

	// UnregisterFromTicksListeningWorkflowResults is the results of the
	// UnregisterFromTicksReception workflow.
	UnregisterFromTicksListeningWorkflowResults struct{}
)
