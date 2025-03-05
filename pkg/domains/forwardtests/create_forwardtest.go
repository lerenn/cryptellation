package forwardtests

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardtestWorkflow creates a new forwardtest and saves it to the database.
func (wf *workflows) CreateForwardtestWorkflow(
	ctx workflow.Context,
	params api.CreateForwardtestWorkflowParams,
) (api.CreateForwardtestWorkflowResults, error) {
	payload := forwardtest.NewForwardtestParams{
		Accounts: params.Accounts,
	}
	if err := payload.Validate(); err != nil {
		return api.CreateForwardtestWorkflowResults{}, err
	}

	// Create new forwardtest and save it to database
	ft := forwardtest.New(payload)
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.CreateForwardtestActivity, db.CreateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardtestWorkflowResults{}, fmt.Errorf("adding forwardtest to db: %w", err)
	}

	return api.CreateForwardtestWorkflowResults{
		ID: ft.ID,
	}, nil
}
