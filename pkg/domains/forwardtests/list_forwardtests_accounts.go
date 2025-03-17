package forwardtests

import (
	"fmt"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"go.temporal.io/sdk/workflow"
)

// ListForwardtestAccountsWorkflow lists the account of a forwardtest.
func (wf *workflows) ListForwardtestAccountsWorkflow(
	ctx workflow.Context,
	params api.ListForwardtestAccountsWorkflowParams,
) (api.ListForwardtestAccountsWorkflowResults, error) {
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return api.ListForwardtestAccountsWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	return api.ListForwardtestAccountsWorkflowResults{
		Accounts: ft.Accounts,
	}, nil
}
