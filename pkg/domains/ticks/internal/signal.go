package internal

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
)

const RegisterToTicksListeningSignalName = "RegisterToTicksListeningSignal"

type (
	RegisterToTicksListeningSignalParams struct {
		CallbackWorkflow api.ListenToTicksCallbackWorkflow
	}
)

const UnregisterFromTicksListeningSignalName = "UnregisterFromTicksListeningSignal"

type (
	UnregisterFromTicksListeningSignalParams struct {
		CallbackWorkflowName string
	}
)

const NewTickReceivedSignalName = "NewTickReceivedSignal"

type (
	NewTickReceivedSignalParams struct {
		Tick tick.Tick
	}
)
