package temporal

import (
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	temporalutils "github.com/lerenn/cryptellation/v1/pkg/temporal"
)

// Callbacks is the struct representing callbacks for ans automation through cryptellation API.
type Callbacks struct {
	OnInitCallback      temporalutils.CallbackWorkflow
	OnNewPricesCallback temporalutils.CallbackWorkflow
	OnExitCallback      temporalutils.CallbackWorkflow
}

// OnInitCallbackWorkflowParams is the parameters of the
// OnInitCallbackWorkflow callback workflow.
type OnInitCallbackWorkflowParams struct {
	RunCtx run.Context
}

// OnNewPricesCallbackWorkflowParams is the parameters of the
// OnNewPricesCallbackWorkflow callback workflow.
type OnNewPricesCallbackWorkflowParams struct {
	Run   run.Context
	Ticks []tick.Tick
}

// OnExitCallbackWorkflowParams is the parameters of the
// OnExitCallbackWorkflow callback workflow.
type OnExitCallbackWorkflowParams struct {
	Run run.Context
}
