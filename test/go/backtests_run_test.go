package test

import (
	"context"
	"time"

	"github.com/google/uuid"
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/go/worker/wfclient"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/run"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.temporal.io/sdk/workflow"
)

type testRobotRun struct {
	Suite *EndToEndSuite

	BacktestID     uuid.UUID
	BacktestParams api.CreateBacktestWorkflowParams

	Cryptellation    wfclient.Client
	OnInitCalls      int
	OnNewPricesCalls int
	OnExitCalls      int
}

func (r *testRobotRun) OnInit(ctx workflow.Context, params api.OnInitCallbackWorkflowParams) error {
	checkBacktestRunContext(r.Suite, params.RunCtx, r.BacktestID)
	r.Suite.Require().WithinDuration(r.BacktestParams.BacktestParameters.StartTime, params.RunCtx.Now, time.Second)

	err := r.Cryptellation.SubscribeToPrice(ctx, wfclient.SubscribeToPriceParams{
		Run:      params.RunCtx,
		Exchange: "binance",
		Pair:     "BTC-USDT",
	})
	r.Suite.Require().NoError(err)

	r.OnInitCalls++
	return err
}

func (r *testRobotRun) OnNewPrices(_ workflow.Context, params api.OnNewPricesCallbackWorkflowParams) error {
	checkBacktestRunContext(r.Suite, params.Run, r.BacktestID)

	// TODO(#72): test order passing in OnNewPrices

	r.OnNewPricesCalls++
	return nil
}

func (r *testRobotRun) OnExit(_ workflow.Context, params api.OnExitCallbackWorkflowParams) error {
	checkBacktestRunContext(r.Suite, params.Run, r.BacktestID)
	r.Suite.Require().WithinDuration(*r.BacktestParams.BacktestParameters.EndTime, params.Run.Now, time.Second)

	r.OnExitCalls++
	return nil
}

func (suite *EndToEndSuite) TestBacktestRun() {
	// WHEN creating a new backtest

	params := api.CreateBacktestWorkflowParams{
		BacktestParameters: backtest.Parameters{
			Accounts: map[string]account.Account{
				"binance": {
					Balances: map[string]float64{
						"BTC": 1,
					},
				},
			},
			StartTime: utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")),
			EndTime:   utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z"))),
		},
	}
	backtest, err := suite.client.NewBacktest(context.Background(), params)

	// THEN no error is returned

	suite.Require().NoError(err)

	// WHEN running the backtest with a robot

	r := &testRobotRun{
		BacktestParams: params,
		BacktestID:     backtest.ID,
		Suite:          suite,
		Cryptellation:  wfclient.NewClient(),
	}
	err = backtest.Run(context.Background(), r)

	// THEN no error is returned

	suite.Require().NoError(err)

	// AND the robot callbacks are called
	suite.Require().Equal(1, r.OnInitCalls)
	suite.Require().Equal(2, r.OnNewPricesCalls)
	suite.Require().Equal(1, r.OnExitCalls)
}

func checkBacktestRunContext(suite *EndToEndSuite, ctx run.Context, backtestID uuid.UUID) {
	suite.Require().Equal(backtestID, ctx.ID)
	suite.Require().Equal(run.ModeBacktest, ctx.Mode)
	suite.Require().NotEmpty(ctx.TaskQueue)
}
