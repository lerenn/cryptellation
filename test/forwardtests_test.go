package test

import (
	"context"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
)

func (suite *EndToEndSuite) TestGetForwardtestStatus() {
	// GIVEN a forwardtest

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN getting the forwardtest status

	status, err := ft.GetStatus(context.Background())
	suite.Require().NoError(err)

	// THEN the status is "running"

	suite.Require().Equal(1000.0, status.Balance)
}

func (suite *EndToEndSuite) TestListForwardtestStatus() {
	// GIVEN 3 forwardtests

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
	}
	ft1, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)
	ft2, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)
	ft3, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN listing the forwardtests

	list, err := suite.client.ListForwardtests(context.Background(), api.ListForwardtestsWorkflowParams{})
	suite.Require().NoError(err)

	// THEN the list contains the forwardtests

	suite.Require().Contains(list, ft1)
	suite.Require().Contains(list, ft2)
	suite.Require().Contains(list, ft3)
}

func (suite *EndToEndSuite) TestCreateOrder() {
	// GIVEN a forwardtest

	params := api.CreateForwardtestWorkflowParams{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000000,
				},
			},
		},
	}
	ft, err := suite.client.NewForwardtest(context.Background(), params)
	suite.Require().NoError(err)

	// WHEN creating an order

	_, err = ft.CreateOrder(context.Background(), order.Order{
		Type:     order.TypeIsMarket,
		Side:     order.SideIsBuy,
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Quantity: 1,
	})
	suite.Require().NoError(err)

	// THEN the balances are in order

	accounts, err := ft.ListAccounts(context.Background())
	suite.Require().NoError(err)
	suite.Require().Equal(1.0, accounts["binance"].Balances["BTC"])
	suite.Require().NotEqual(1000000.0, accounts["binance"].Balances["USDT"])
}
