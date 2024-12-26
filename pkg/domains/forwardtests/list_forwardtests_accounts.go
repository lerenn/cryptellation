package forwardtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/workflow"
)

// ListForwardTestAccountsWorkflow lists the account of a forwardtest.
func (wf *workflows) ListForwardTestAccountsWorkflow(
	ctx workflow.Context,
	params api.ListForwardTestAccountsWorkflowParams,
) (api.ListForwardTestAccountsWorkflowResults, error) {
	ft, err := wf.readForwardTestFromDB(ctx, params.ForwardTestID)
	if err != nil {
		return api.ListForwardTestAccountsWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	return api.ListForwardTestAccountsWorkflowResults{
		Accounts: ft.Accounts,
	}, nil
}
