package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

func (suite *EndToEndSuite) TestBacktestGet() {
	// GIVEN a backtest

	params := api.CreateBacktestWorkflowParams{
		BacktestParameters: backtest.Parameters{
			Accounts: map[string]account.Account{
				"binance": {
					Balances: map[string]float64{
						"BTC": 1,
					},
				},
			},
		},
	}
	bt, err := suite.client.NewBacktest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN getting the backtest

	res, err := suite.client.GetBacktest(context.Background(), api.GetBacktestWorkflowParams{
		BacktestID: bt.ID,
	})
	suite.Require().NoError(err)

	// THEN the response is the backtest

	suite.Require().Equal(bt, res)
}

func (suite *EndToEndSuite) TestBacktestList() {
	// GIVEN 3 backtests

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
	bt1, err := suite.client.NewBacktest(context.Background(), params)
	suite.Require().NoError(err)
	bt2, err := suite.client.NewBacktest(context.Background(), params)
	suite.Require().NoError(err)
	bt3, err := suite.client.NewBacktest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN listing backtests

	res, err := suite.client.ListBacktests(context.Background(), api.ListBacktestsWorkflowParams{})
	suite.Require().NoError(err)

	// THEN the response contains the 3 backtests

	suite.Require().Contains(res, bt1)
	suite.Require().Contains(res, bt2)
	suite.Require().Contains(res, bt3)
}
