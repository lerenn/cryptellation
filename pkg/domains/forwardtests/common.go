package forwardtests

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) readForwardtestFromDB(ctx workflow.Context, id uuid.UUID) (forwardtest.Forwardtest, error) {
	var readRes db.ReadForwardtestActivityResult
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadForwardtestActivity, db.ReadForwardtestActivityParams{
			ID: id,
		}).Get(ctx, &readRes)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	return readRes.Forwardtest, nil
}
