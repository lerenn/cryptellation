package backtests

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/domains/backtests/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	temporalutils "github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) RunBacktestWorkflow(
	ctx workflow.Context,
	params api.RunBacktestWorkflowParams,
) (api.RunBacktestWorkflowResults, error) {
	logger := workflow.GetLogger(ctx)

	// Getting backtest
	var dbRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
		ID: params.BacktestID,
	}).Get(ctx, &dbRes)
	if err != nil {
		return api.RunBacktestWorkflowResults{}, fmt.Errorf("reading backtest from db: %w", err)
	}
	bt := dbRes.Backtest

	// Init the backtest from client side
	if err := initBacktestCallback(ctx, params.Callbacks.OnInitCallback); err != nil {
		return api.RunBacktestWorkflowResults{}, fmt.Errorf("initializing backtest from client side: %w", err)
	}

	// Loop until backtest is finished
	for finished := false; !finished; {
		// Get prices
		prices, err := wf.readActualPrices(ctx, bt)
		if err != nil {
			return api.RunBacktestWorkflowResults{}, fmt.Errorf("cannot read actual prices: %w", err)
		}
		if len(prices) == 0 {
			logger.Warn("No price detected",
				"time", bt.CurrentCandlestick.Time)
			bt.SetCurrentTime(bt.Parameters.EndTime)
			break
		} else if !prices[0].Time.Equal(bt.CurrentCandlestick.Time) {
			logger.Warn("No price between current time and first event retrieved",
				"current_time", bt.CurrentCandlestick.Time,
				"first_event_time", prices[0].Time)
			bt.SetCurrentTime(prices[0].Time)
		}

		// Execute backtest with these prices
		if err := executeBacktest(ctx, params.Callbacks.OnNewPricesCallback, prices); err != nil {
			return api.RunBacktestWorkflowResults{}, fmt.Errorf("cannot execute backtest: %w", err)
		}

		// Advance backtest
		finished, bt, err = wf.advanceBacktest(ctx, bt.ID)
		if err != nil {
			return api.RunBacktestWorkflowResults{}, fmt.Errorf("cannot advance backtest: %w", err)
		}
	}

	// Exit the backtest from client side
	if err := exitBacktestCallback(ctx, params.Callbacks.OnExitCallback); err != nil {
		return api.RunBacktestWorkflowResults{}, fmt.Errorf("exit backtest from client side: %w", err)
	}

	return api.RunBacktestWorkflowResults{}, nil
}

func (wf *workflows) advanceBacktest(ctx workflow.Context, id uuid.UUID) (bool, backtest.Backtest, error) {
	logger := workflow.GetLogger(ctx)

	// Read backtest
	var readRes db.ReadBacktestActivityResults
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.ReadBacktestActivity, db.ReadBacktestActivityParams{
		ID: id,
	}).Get(ctx, &readRes)
	if err != nil {
		return false, backtest.Backtest{}, fmt.Errorf("read backtest from db: %w", err)
	}
	bt := readRes.Backtest

	// Advance backtest
	finished, err := bt.Advance()
	if err != nil {
		return false, backtest.Backtest{}, fmt.Errorf("cannot advance backtest: %w", err)
	}
	logger.Info("Advancing backtest",
		"id", bt.ID.String(),
		"current_time", bt.CurrentTime())

	// Save backtest
	var writeRes db.UpdateBacktestActivityResults
	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: activities.DBInteractionDefaultTimeout,
	}), wf.db.UpdateBacktestActivity, db.UpdateBacktestActivityParams{
		Backtest: bt,
	}).Get(ctx, &writeRes)
	if err != nil {
		return false, backtest.Backtest{}, fmt.Errorf("save backtest to db: %w", err)
	}

	return finished, bt, nil
}

func (wf *workflows) readActualPrices(ctx workflow.Context, bt backtest.Backtest) ([]tick.Tick, error) {
	logger := workflow.GetLogger(ctx)

	// Run for all prices subscriptions
	// TODO: parallelize
	prices := make([]tick.Tick, 0, len(bt.PricesSubscriptions))
	for _, sub := range bt.PricesSubscriptions {
		// Get exchange info
		var result api.ListCandlesticksWorkflowResults
		if err := workflow.ExecuteChildWorkflow(ctx, api.ListCandlesticksWorkflowName, api.ListCandlesticksWorkflowParams{
			Exchange: sub.Exchange,
			Pair:     sub.Pair,
			Period:   bt.Parameters.PricePeriod,
			Start:    &bt.CurrentCandlestick.Time,
			End:      &bt.Parameters.EndTime,
			Limit:    1,
		}).Get(ctx, &result); err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		// Get the first candlestick if possible
		t, cs, exists := result.List.Data.First()
		if !exists {
			continue
		}

		// Create tick from candlesticks
		p := tick.FromCandlestick(sub.Exchange, sub.Pair, bt.CurrentCandlestick.Price, t, cs)
		prices = append(prices, p)
	}

	// Only keep the earliest same time ticks for time consistency
	_, prices = tick.OnlyKeepEarliestSameTime(prices, bt.Parameters.EndTime)
	logger.Info("Gotten ticks on backtest",
		"quantity", len(prices),
		"backtest_id", bt.ID.String())
	return prices, nil
}

func executeBacktest(
	ctx workflow.Context,
	callback temporalutils.CallbackWorkflow,
	prices []tick.Tick,
) error {
	// Options
	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                callback.TaskQueueName, // Execute in the client queue
		WorkflowExecutionTimeout: time.Second * 30,       // Timeout if the child workflow does not complete
	}

	// Check if the timeout is set
	if callback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = callback.ExecutionTimeout
	}

	// Execute backtest
	err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, opts),
		callback.Name, api.OnNewPricesCallbackWorkflowParams{
			Ticks: prices,
		}).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func exitBacktestCallback(
	ctx workflow.Context,
	onExitCallback temporalutils.CallbackWorkflow,
) error {
	// Options
	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                onExitCallback.TaskQueueName, // Execute in the client queue
		WorkflowExecutionTimeout: time.Second * 30,             // Timeout if the child workflow does not complete
	}

	// Check if the timeout is set
	if onExitCallback.ExecutionTimeout > 0 {
		opts.WorkflowExecutionTimeout = onExitCallback.ExecutionTimeout
	}

	// Run a new child workflow
	ctx = workflow.WithChildOptions(ctx, opts)
	if err := workflow.ExecuteChildWorkflow(
		ctx, onExitCallback.Name, api.OnExitCallbackWorkflowParams{},
	).Get(ctx, nil); err != nil {
		return fmt.Errorf("starting new onExitCallback child workflow: %w", err)
	}

	return nil
}
