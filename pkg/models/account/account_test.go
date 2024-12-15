package account

import (
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/models/order"
	"github.com/stretchr/testify/suite"
)

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

type AccountSuite struct {
	suite.Suite
}

func (suite *AccountSuite) TestValidate() {
	acc := Account{
		Balances: map[string]float64{
			"USDC": 1000,
			"ETH":  1,
		},
	}

	suite.Require().NoError(acc.Validate())
}

func (suite *AccountSuite) TestValidateWithInvalidName() {
	acc := Account{
		Balances: map[string]float64{
			"":    1000,
			"ETH": 1,
		},
	}

	suite.Require().Error(acc.Validate())
}

func (suite *AccountSuite) TestValidateWithInvalidValue() {
	acc := Account{
		Balances: map[string]float64{
			"USDC": 1000,
			"ETH":  -1,
		},
	}

	suite.Require().Error(acc.Validate())
}

func (suite *AccountSuite) TestApplyOrder() {
	cases := []struct {
		before Account
		price  float64
		order  order.Order
		after  Account
	}{
		// Buy
		{
			before: Account{
				Balances: map[string]float64{"USDC": 100, "ETH": 1},
			},
			price: 100,
			order: order.Order{
				Pair:     "ETH-USDC",
				Quantity: 1,
				Side:     order.SideIsBuy,
			},
			after: Account{
				Balances: map[string]float64{"USDC": 0, "ETH": 2},
			},
		},
		// Sell
		{
			before: Account{
				Balances: map[string]float64{"USDC": 100, "ETH": 1},
			},
			price: 100,
			order: order.Order{
				Pair:     "ETH-USDC",
				Quantity: 1,
				Side:     order.SideIsSell,
			},
			after: Account{
				Balances: map[string]float64{"USDC": 200, "ETH": 0},
			},
		},
	}

	for _, c := range cases {
		suite.Require().NoError(c.before.ApplyOrder(c.price, c.order))
		suite.Require().Equal(c.after, c.before)
	}
}
