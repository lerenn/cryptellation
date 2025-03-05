package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/stretchr/testify/suite"
)

// ForwardtestSuite is the suite test for forwardtest db activities.
type ForwardtestSuite struct {
	suite.Suite
	DB DB
}

// TestCreateReadForwardtestActivities tests the create and read operations.
func (suite *ForwardtestSuite) TestCreateReadForwardtestActivities() {
	ft := forwardtest.Forwardtest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateForwardtestActivity(context.Background(), CreateForwardtestActivityParams{
		Forwardtest: ft,
	})
	suite.Require().NoError(err)
	rp, err := suite.DB.ReadForwardtestActivity(context.Background(), ReadForwardtestActivityParams{
		ID: ft.ID,
	})
	suite.Require().NoError(err, ft.ID.String())

	suite.Require().Equal(ft.ID, rp.Forwardtest.ID)
	suite.Require().Len(rp.Forwardtest.Accounts, 1)
	suite.Require().Len(rp.Forwardtest.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(
		ft.Accounts["exchange"].Balances["DAI"],
		rp.Forwardtest.Accounts["exchange"].Balances["DAI"])
}

// TestListForwardtestsActivity tests the list operation.
func (suite *ForwardtestSuite) TestListForwardtestsActivity() {
	ft1 := forwardtest.Forwardtest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateForwardtestActivity(context.Background(), CreateForwardtestActivityParams{
		Forwardtest: ft1,
	})
	suite.Require().NoError(err)
	ft2 := forwardtest.Forwardtest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1500,
				},
			},
		},
	}
	_, err = suite.DB.CreateForwardtestActivity(context.Background(), CreateForwardtestActivityParams{
		Forwardtest: ft2,
	})
	suite.Require().NoError(err)

	rp, err := suite.DB.ListForwardtestsActivity(context.Background(), ListForwardtestsActivityParams{})
	suite.Require().NoError(err)

	suite.Require().Len(rp.Forwardtests, 2)
	suite.Require().Equal(rp.Forwardtests[0].ID, ft2.ID) // Last created first
	suite.Require().Equal(rp.Forwardtests[1].ID, ft1.ID)
}

// TestUpdateForwardtestActivity tests the update operation.
func (suite *ForwardtestSuite) TestUpdateForwardtestActivity() {
	// Create forwardtest
	ft1 := forwardtest.Forwardtest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateForwardtestActivity(context.Background(), CreateForwardtestActivityParams{
		Forwardtest: ft1,
	})
	suite.Require().NoError(err)
	rp1, err := suite.DB.ReadForwardtestActivity(context.Background(), ReadForwardtestActivityParams{
		ID: ft1.ID,
	})
	suite.Require().NoError(err)

	// Wait for 1 millisecond
	time.Sleep(time.Millisecond)

	// Update forwardtest
	ft2 := forwardtest.Forwardtest{
		ID: ft1.ID,
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	_, err = suite.DB.UpdateForwardtestActivity(context.Background(), UpdateForwardtestActivityParams{
		Forwardtest: ft2,
	})
	suite.Require().NoError(err)
	rp2, err := suite.DB.ReadForwardtestActivity(context.Background(), ReadForwardtestActivityParams{
		ID: ft1.ID,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(ft1.ID, rp2.Forwardtest.ID)
	suite.Require().True(
		rp2.Forwardtest.UpdatedAt.After(rp1.Forwardtest.UpdatedAt),
		rp2.Forwardtest.UpdatedAt.String()+" should be after "+rp1.Forwardtest.UpdatedAt.String())
	suite.Require().Equal(ft2.ID, rp2.Forwardtest.ID)
	suite.Require().Len(rp2.Forwardtest.Accounts, 1)
	suite.Require().Len(rp2.Forwardtest.Accounts["exchange2"].Balances, 1)
	suite.Require().Equal(
		ft2.Accounts["exchange2"].Balances["USDC"],
		rp2.Forwardtest.Accounts["exchange2"].Balances["USDC"])
}

// TestDeleteForwardtestActivity tests the delete operation.
func (suite *ForwardtestSuite) TestDeleteForwardtestActivity() {
	ft := forwardtest.Forwardtest{
		ID: uuid.New(),
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateForwardtestActivity(context.Background(), CreateForwardtestActivityParams{
		Forwardtest: ft,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.DeleteForwardtestActivity(context.Background(), DeleteForwardtestActivityParams{
		ID: ft.ID,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.ReadForwardtestActivity(context.Background(), ReadForwardtestActivityParams{
		ID: ft.ID,
	})
	suite.Error(err)
}
