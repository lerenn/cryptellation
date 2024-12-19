package bot

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

// Bot is the interface for a trading robot executed on Cryptellation.
type Bot interface {
	OnInit(ctx workflow.Context, params api.OnInitCallbackWorkflowParams) error
	OnNewPrices(ctx workflow.Context, params api.OnNewPricesCallbackWorkflowParams) error
	OnExit(ctx workflow.Context, params api.OnExitCallbackWorkflowParams) error
}
