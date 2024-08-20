package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"

	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ForwardTestSuite struct {
	suite.Suite
	DB Port
}

func (suite *ForwardTestSuite) TestCreateRead() {
	ft := forwardtest.ForwardTest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateForwardTest(context.TODO(), ft))
	rp, err := suite.DB.ReadForwardTest(context.TODO(), ft.ID)
	suite.Require().NoError(err, ft.ID.String())

	suite.Require().Equal(ft.ID, rp.ID)
	suite.Require().Len(rp.Accounts, 1)
	suite.Require().Len(rp.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(ft.Accounts["exchange"].Balances["DAI"], rp.Accounts["exchange"].Balances["DAI"])
}

func (suite *ForwardTestSuite) TestListSuite() {
	ft1 := forwardtest.ForwardTest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateForwardTest(context.TODO(), ft1))
	ft2 := forwardtest.ForwardTest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1500,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateForwardTest(context.TODO(), ft2))

	rp, err := suite.DB.ListForwardTests(context.TODO(), ListFilters{})
	suite.Require().NoError(err)

	suite.Require().Len(rp, 2)
	suite.Require().Equal(rp[0].ID, ft2.ID) // Last created first
	suite.Require().Equal(rp[1].ID, ft1.ID)
}

func (suite *ForwardTestSuite) TestUpdate() {
	// Create forward test
	ft1 := forwardtest.ForwardTest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateForwardTest(context.TODO(), ft1))
	rp1, err := suite.DB.ReadForwardTest(context.TODO(), ft1.ID)
	suite.Require().NoError(err)

	// Wait for 1 millisecond
	time.Sleep(time.Millisecond)

	// Update forward test
	ft2 := forwardtest.ForwardTest{
		ID: ft1.ID,
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.UpdateForwardTest(context.TODO(), ft2))
	rp2, err := suite.DB.ReadForwardTest(context.TODO(), ft1.ID)
	suite.Require().NoError(err)

	suite.Require().Equal(ft1.ID, rp2.ID)
	suite.Require().True(rp2.UpdatedAt.After(rp1.UpdatedAt), rp2.UpdatedAt.String()+" should be after "+rp1.UpdatedAt.String())
	suite.Require().Equal(ft2.ID, rp2.ID)
	suite.Require().Len(rp2.Accounts, 1)
	suite.Require().Len(rp2.Accounts["exchange2"].Balances, 1)
	suite.Require().Equal(ft2.Accounts["exchange2"].Balances["USDC"], rp2.Accounts["exchange2"].Balances["USDC"])
}

func (suite *ForwardTestSuite) TestDelete() {
	ft := forwardtest.ForwardTest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateForwardTest(context.TODO(), ft))
	suite.Require().NoError(suite.DB.DeleteForwardTest(context.TODO(), ft.ID))
	_, err := suite.DB.ReadForwardTest(context.TODO(), ft.ID)
	suite.Error(err)
}
