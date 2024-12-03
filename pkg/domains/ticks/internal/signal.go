package internal

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
)

// RegisterToTicksListeningSignalName is the name of the signal to send when
// registering to ticks listening through the corresponding workflow.
const RegisterToTicksListeningSignalName = "RegisterToTicksListeningSignal"

type (
	// RegisterToTicksListeningSignalParams is the parameters of the RegisterToTicksListeningSignal.
	RegisterToTicksListeningSignalParams struct {
		CallbackWorkflow api.ListenToTicksCallbackWorkflow
	}
)

// UnregisterFromTicksListeningSignalName is the name of the signal to send when
// unregistering from ticks listening through the corresponding workflow.
const UnregisterFromTicksListeningSignalName = "UnregisterFromTicksListeningSignal"

type (
	// UnregisterFromTicksListeningSignalParams is the parameters of the UnregisterFromTicksListeningSignal.
	UnregisterFromTicksListeningSignalParams struct {
		CallbackWorkflowName string
	}
)

// NewTickReceivedSignalName is the name of the signal to send when a new tick
// is received from an exchange.
const NewTickReceivedSignalName = "NewTickReceivedSignal"

type (
	// NewTickReceivedSignalParams is the parameters of the NewTickReceivedSignal.
	NewTickReceivedSignalParams struct {
		Tick tick.Tick
	}
)
