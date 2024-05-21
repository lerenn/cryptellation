package test

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
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
