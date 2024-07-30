package test

import (
	"context"

	common "cryptellation/pkg/client"
	"cryptellation/pkg/models/account"
	"cryptellation/pkg/models/order"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
)

func (suite *EndToEndSuite) TestCreateForwardTest() {
	// Create forwardtest
	id, err := suite.client.CreateForwardTest(context.Background(), forwardtest.NewPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"BTC": 1,
				},
			},
		},
	})
	suite.Require().NoError(err)
	suite.Require().NotEqual(uuid.Nil, id)
}

func (suite *EndToEndSuite) TestCreateOrder() {
	// Create forwardtest
	id, err := suite.client.CreateForwardTest(context.Background(), forwardtest.NewPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000000,
				},
			},
		},
	})
	suite.Require().NoError(err)

	// Create order
	err = suite.client.CreateOrder(context.Background(), common.OrderCreationPayload{
		RunID:    id,
		Type:     order.TypeIsMarket,
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Side:     order.SideIsBuy,
		Quantity: 1,
	})
	suite.Require().NoError(err)

	// Check balances
	accounts, err := suite.client.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(1.0, accounts["binance"].Balances["BTC"])
	suite.Require().NotEqual(1000000, accounts["binance"].Balances["USDT"])
}
