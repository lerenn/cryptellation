package forwardtests

import (
	"fmt"

	"github.com/google/uuid"
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardtestOrderWorkflow creates a new forwardtest order and saves it to the database.
func (wf *workflows) CreateForwardtestOrderWorkflow(
	ctx workflow.Context,
	params api.CreateForwardtestOrderWorkflowParams,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)

	if params.Order.ID == uuid.Nil {
		params.Order.ID = uuid.New()
	}

	logger.Debug("Creating order on forwardtest",
		"order", params.Order,
		"forwardtest_id", params.ForwardtestID.String())

	// Read forwardtest from database
	ft, err := wf.readForwardtestFromDB(ctx, params.ForwardtestID)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{},
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
		return api.CreateForwardtestOrderWorkflowResults{},
			fmt.Errorf("could not get candlesticks from service: %w", err)
	}

	cs, ok := csRes.List.First()
	if !ok {
		return api.CreateForwardtestOrderWorkflowResults{}, fmt.Errorf("no data for order validation")
	}

	logger.Info("Adding order to forwardtest",
		"order", params.Order,
		"forwardtest", params.ForwardtestID.String())
	if err := ft.AddOrder(params.Order, cs); err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	// Save forwardtest to database
	err = workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, db.DefaultActivityOptions()),
		wf.db.UpdateForwardtestActivity, db.UpdateForwardtestActivityParams{
			Forwardtest: ft,
		}).Get(ctx, nil)
	if err != nil {
		return api.CreateForwardtestOrderWorkflowResults{}, err
	}

	return api.CreateForwardtestOrderWorkflowResults{}, nil
}
