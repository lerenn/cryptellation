package test

import (
	"context"
	"time"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/order"

	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

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

func (suite *EndToEndSuite) TestListForwardTests() {
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

	// List forwardtests
	tests, err := suite.client.ListForwardTests(context.Background())
	suite.Require().NoError(err)
	suite.Require().True(len(tests) > 0)

	// Get last updated forwardtest
	var last forwardtest.ForwardTest
	for _, t := range tests {
		if t.UpdatedAt.After(last.UpdatedAt) {
			last = t
		}
	}

	// Check last forwardtest
	suite.Require().Equal(id, last.ID)
	suite.Require().WithinDuration(time.Now(), last.UpdatedAt, 5*time.Second)
}
