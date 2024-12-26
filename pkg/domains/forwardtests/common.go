package forwardtests

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) readForwardTestFromDB(ctx workflow.Context, id uuid.UUID) (forwardtest.ForwardTest, error) {
	var readRes db.ReadForwardTestActivityResult
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.ReadForwardTestActivity, db.ReadForwardTestActivityParams{
			ID: id,
		}).Get(ctx, &readRes)
	if err != nil {
		return forwardtest.ForwardTest{}, err
	}

	return readRes.ForwardTest, nil
}
