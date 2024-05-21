package test

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/order"
	client "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

func (suite *EndToEndSuite) CreateForwardTest() {
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

func (suite *EndToEndSuite) CreateOrder() {
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
	err = suite.client.CreateOrder(context.Background(), client.OrderCreationPayload{
		ForwardTestID: id,
		Type:          order.TypeIsMarket,
		Exchange:      "binance",
		Pair:          "BTC-USDT",
		Side:          order.SideIsBuy,
		Quantity:      1,
	})
	suite.Require().NoError(err)

	// Check balances
	// TODO
}
