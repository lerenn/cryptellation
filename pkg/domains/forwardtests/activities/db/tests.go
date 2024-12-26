package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/stretchr/testify/suite"
)

// ForwardTestSuite is the suite test for forwardtest db activities.
type ForwardTestSuite struct {
	suite.Suite
	DB DB
}

// TestCreateReadForwardTestActivities tests the create and read operations.
func (suite *ForwardTestSuite) TestCreateReadForwardTestActivities() {
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
	_, err := suite.DB.CreateForwardTestActivity(context.Background(), CreateForwardTestActivityParams{
		ForwardTest: ft,
	})
	suite.Require().NoError(err)
	rp, err := suite.DB.ReadForwardTestActivity(context.Background(), ReadForwardTestActivityParams{
		ID: ft.ID,
	})
	suite.Require().NoError(err, ft.ID.String())

	suite.Require().Equal(ft.ID, rp.ForwardTest.ID)
	suite.Require().Len(rp.ForwardTest.Accounts, 1)
	suite.Require().Len(rp.ForwardTest.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(
		ft.Accounts["exchange"].Balances["DAI"],
		rp.ForwardTest.Accounts["exchange"].Balances["DAI"])
}

// TestListForwardTestsActivity tests the list operation.
func (suite *ForwardTestSuite) TestListForwardTestsActivity() {
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
	_, err := suite.DB.CreateForwardTestActivity(context.Background(), CreateForwardTestActivityParams{
		ForwardTest: ft1,
	})
	suite.Require().NoError(err)
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
	_, err = suite.DB.CreateForwardTestActivity(context.Background(), CreateForwardTestActivityParams{
		ForwardTest: ft2,
	})
	suite.Require().NoError(err)

	rp, err := suite.DB.ListForwardTestsActivity(context.Background(), ListForwardTestsActivityParams{})
	suite.Require().NoError(err)

	suite.Require().Len(rp.ForwardTests, 2)
	suite.Require().Equal(rp.ForwardTests[0].ID, ft2.ID) // Last created first
	suite.Require().Equal(rp.ForwardTests[1].ID, ft1.ID)
}

// TestUpdateForwardTestActivity tests the update operation.
func (suite *ForwardTestSuite) TestUpdateForwardTestActivity() {
	// Create forwardtest
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
	_, err := suite.DB.CreateForwardTestActivity(context.Background(), CreateForwardTestActivityParams{
		ForwardTest: ft1,
	})
	suite.Require().NoError(err)
	rp1, err := suite.DB.ReadForwardTestActivity(context.Background(), ReadForwardTestActivityParams{
		ID: ft1.ID,
	})
	suite.Require().NoError(err)

	// Wait for 1 millisecond
	time.Sleep(time.Millisecond)

	// Update forwardtest
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
	_, err = suite.DB.UpdateForwardTestActivity(context.Background(), UpdateForwardTestActivityParams{
		ForwardTest: ft2,
	})
	suite.Require().NoError(err)
	rp2, err := suite.DB.ReadForwardTestActivity(context.Background(), ReadForwardTestActivityParams{
		ID: ft1.ID,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(ft1.ID, rp2.ForwardTest.ID)
	suite.Require().True(
		rp2.ForwardTest.UpdatedAt.After(rp1.ForwardTest.UpdatedAt),
		rp2.ForwardTest.UpdatedAt.String()+" should be after "+rp1.ForwardTest.UpdatedAt.String())
	suite.Require().Equal(ft2.ID, rp2.ForwardTest.ID)
	suite.Require().Len(rp2.ForwardTest.Accounts, 1)
	suite.Require().Len(rp2.ForwardTest.Accounts["exchange2"].Balances, 1)
	suite.Require().Equal(
		ft2.Accounts["exchange2"].Balances["USDC"],
		rp2.ForwardTest.Accounts["exchange2"].Balances["USDC"])
}

// TestDeleteForwardTestActivity tests the delete operation.
func (suite *ForwardTestSuite) TestDeleteForwardTestActivity() {
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
	_, err := suite.DB.CreateForwardTestActivity(context.Background(), CreateForwardTestActivityParams{
		ForwardTest: ft,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.DeleteForwardTestActivity(context.Background(), DeleteForwardTestActivityParams{
		ID: ft.ID,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.ReadForwardTestActivity(context.Background(), ReadForwardTestActivityParams{
		ID: ft.ID,
	})
	suite.Error(err)
}
