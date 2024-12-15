package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.temporal.io/sdk/workflow"
)

type testRobotCallbacks struct {
	OnInitCalls      int
	OnNewPricesCalls int
	OnExitCalls      int
}

func (r *testRobotCallbacks) OnInit(_ workflow.Context, _ api.OnInitCallbackWorkflowParams) error {
	r.OnInitCalls++
	return nil
}

func (r *testRobotCallbacks) OnNewPrices(_ workflow.Context, _ api.OnNewPricesCallbackWorkflowParams) error {
	r.OnNewPricesCalls++
	return nil
}

func (r *testRobotCallbacks) OnExit(_ workflow.Context, _ api.OnExitCallbackWorkflowParams) error {
	r.OnExitCalls++
	return nil
}

func (suite *EndToEndSuite) TestBacktestCallbacks() {
	// WHEN creating a new backtest

	backtest, err := suite.client.NewBacktest(context.Background(), api.CreateBacktestWorkflowParams{
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
	})

	// THEN no error is returned

	suite.Require().NoError(err)

	// WHEN running the backtest with a robot

	r := &testRobotCallbacks{}
	err = backtest.Run(context.Background(), r)

	// THEN no error is returned

	suite.Require().NoError(err)

	// AND the robot callbacks are called
	suite.Require().Equal(1, r.OnInitCalls)
	suite.Require().Equal(2, r.OnNewPricesCalls)
	suite.Require().Equal(1, r.OnExitCalls)
}
