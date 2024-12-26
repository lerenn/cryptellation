package forwardtests

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardTestOrderWorkflow creates a new forwardtest order and saves it to the database.
func (wf *workflows) CreateForwardTestOrderWorkflow(
	ctx workflow.Context,
	params api.CreateForwardTestOrderWorkflowParams,
) (api.CreateForwardTestOrderWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)

	if params.Order.ID == uuid.Nil {
		params.Order.ID = uuid.New()
	}

	logger.Debug("Creating order on forwardtest",
		"order", params.Order,
		"forwardtest_id", params.ForwardTestID.String())

	// Read forwardtest from database
	ft, err := wf.readForwardTestFromDB(ctx, params.ForwardTestID)
	if err != nil {
		return api.CreateForwardTestOrderWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	// Get candlestick for order validation
	now := workflow.Now(ctx)
	csRes, err := wf.cryptellation.ListCandlesticks(ctx, api.ListCandlesticksWorkflowParams{
		Exchange: params.Order.Exchange,
		Pair:     params.Order.Pair,
		Period:   period.M1,
		Start:    &now,
		End:      &now,
		Limit:    1,
	}, nil)
	if err != nil {
		return api.CreateForwardTestOrderWorkflowResults{},
			fmt.Errorf("could not get candlesticks from service: %w", err)
	}

	cs, ok := csRes.List.First()
	if !ok {
		return api.CreateForwardTestOrderWorkflowResults{}, fmt.Errorf("no data for order validation")
	}

	logger.Info("Adding order to forwardtest",
		"order", params.Order,
		"forwardtest", params.ForwardTestID.String())
	if err := ft.AddOrder(params.Order, cs); err != nil {
		return api.CreateForwardTestOrderWorkflowResults{}, err
	}

	// Save forwardtest to database
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.UpdateForwardTestActivity, db.UpdateForwardTestActivityParams{
			ForwardTest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardTestOrderWorkflowResults{}, err
	}

	return api.CreateForwardTestOrderWorkflowResults{}, nil
}
