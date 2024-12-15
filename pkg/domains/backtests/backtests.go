package backtests

import (
	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Backtests is the interface for the backtests domain.
type Backtests interface {
	Register(w worker.Worker)

	// Backtests

	CreateBacktestWorkflow(
		ctx workflow.Context,
		params api.CreateBacktestWorkflowParams,
	) (api.CreateBacktestWorkflowResults, error)
	GetBacktestWorkflow(
		ctx workflow.Context,
		params api.GetBacktestWorkflowParams,
	) (api.GetBacktestWorkflowResults, error)
	ListBacktestsWorkflow(
		ctx workflow.Context,
		params api.ListBacktestsWorkflowParams,
	) (api.ListBacktestsWorkflowResults, error)
	RunBacktestWorkflow(
		ctx workflow.Context,
		params api.RunBacktestWorkflowParams,
	) (api.RunBacktestWorkflowResults, error)

	// Backtests Accounts

	GetBacktestAccountsWorkflow(
		ctx workflow.Context,
		params api.GetBacktestAccountsWorkflowParams,
	) (api.GetBacktestAccountsWorkflowResults, error)

	// Backtests Orders

	CreateBacktestOrderWorkflow(
		ctx workflow.Context,
		params api.CreateBacktestOrderWorkflowParams,
	) (api.CreateBacktestOrderWorkflowResults, error)

	GetBacktestOrdersWorkflow(
		ctx workflow.Context,
		params api.GetBacktestOrdersWorkflowParams,
	) (api.GetBacktestOrdersWorkflowResults, error)
}
