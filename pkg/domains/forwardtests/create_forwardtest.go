package forwardtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardTestWorkflow creates a new forwardtest and saves it to the database.
func (wf *workflows) CreateForwardTestWorkflow(
	ctx workflow.Context,
	params api.CreateForwardTestWorkflowParams,
) (api.CreateForwardTestWorkflowResults, error) {
	payload := forwardtest.NewForwardTestParams{
		Accounts: params.Accounts,
	}
	if err := payload.Validate(); err != nil {
		return api.CreateForwardTestWorkflowResults{}, err
	}

	// Create new forwardtest and save it to database
	ft := forwardtest.New(payload)
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.CreateForwardTestActivity, db.CreateForwardTestActivityParams{
			ForwardTest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardTestWorkflowResults{}, fmt.Errorf("adding forwardtest to db: %w", err)
	}

	return api.CreateForwardTestWorkflowResults{
		ID: ft.ID,
	}, nil
}
