package account

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
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

func (suite *AccountSuite) TestFromProtoBuf() {
	a := FromProtoBuf(&proto.Account{
		Assets: map[string]float64{
			"asset": 32,
		},
	})

	suite.Require().Len(a.Balances, 1)
	suite.Require().Equal(float64(32), a.Balances["asset"])
}

func (suite *AccountSuite) TestToProtoBuf() {
	a := Account{
		Balances: map[string]float64{
			"asset": 32,
		},
	}

	pb := a.ToProtoBuf()
	suite.Require().Len(pb.Assets, 1)
	suite.Require().Equal(float64(32), pb.Assets["asset"])
}
