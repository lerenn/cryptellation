package internal

import (
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
)

const RegisterToTicksReceptionSignalName = "RegisterToTicksReceptionSignal"

type (
	RegisterToTicksReceptionSignalParams struct {
		CallbackWorkflow api.RegisterForTicksCallbackWorkflow
	}
)

const NewTickReceivedSignalName = "NewTickReceivedSignal"

type (
	NewTickReceivedSignalParams struct {
		Tick tick.Tick
	}
)
