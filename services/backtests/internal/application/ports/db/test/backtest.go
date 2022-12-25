package test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/stretchr/testify/suite"
)

type BacktestSuite struct {
	suite.Suite
	DB db.Adapter
}

func (suite *BacktestSuite) TestCreateRead() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	suite.NoError(suite.DB.CreateBacktest(context.TODO(), &bt))
	rp, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Assert().NoError(err)

	suite.Assert().Equal(bt.ID, rp.ID)
	suite.Assert().Len(rp.Accounts, 1)
	suite.Assert().Len(rp.Accounts["exchange"].Balances, 1)
	suite.Assert().Equal(bt.Accounts["exchange"].Balances["DAI"], rp.Accounts["exchange"].Balances["DAI"])
}

func (suite *BacktestSuite) TestUpdate() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.DB.CreateBacktest(context.TODO(), &bt))
	bt2 := backtest.Backtest{
		ID: bt.ID,
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	// Should be changes here
	suite.NoError(suite.DB.UpdateBacktest(context.TODO(), bt2))
	rp, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Assert().NoError(err)

	suite.Equal(bt.ID, rp.ID)
	suite.Equal(bt2.ID, rp.ID)
	suite.Assert().Len(rp.Accounts, 1)
	suite.Assert().Len(rp.Accounts["exchange2"].Balances, 1)
	suite.Assert().Equal(bt2.Accounts["exchange2"].Balances["USDC"], rp.Accounts["exchange2"].Balances["USDC"])
}

func (suite *BacktestSuite) TestDelete() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.DB.CreateBacktest(context.TODO(), &bt))
	suite.NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
	_, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Error(err)
}

func (suite *BacktestSuite) TestDeleteInexistant() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.NoError(suite.DB.CreateBacktest(context.TODO(), &bt))
	suite.NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
	suite.NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
}

func (suite *BacktestSuite) TestLock() {
	bt := backtest.Backtest{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), &bt))

	// Check that lock works even with panic
	suite.Require().NoError(suite.DB.LockedBacktest(bt.ID, func() error {
		panic(errors.New("PANIC"))
	}))

	for i := 0; i < 10; i++ {
		suite.Require().NoError(suite.DB.LockedBacktest(bt.ID, func() error {
			return nil
		}), fmt.Sprintf("Lock/Unlock attempt #%d", i+1))
	}

	go func() {
		err := suite.DB.LockedBacktest(bt.ID, func() error {
			bt.Accounts["exchange"].Balances["ETH"] = 2000
			time.Sleep(200 * time.Millisecond)
			suite.Require().NoError(suite.DB.UpdateBacktest(context.TODO(), bt))
			return nil
		})
		suite.Require().NoError(err)
	}()
	time.Sleep(time.Millisecond)

	suite.Require().NoError(suite.DB.LockedBacktest(bt.ID, func() error {
		recvBt, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
		suite.Require().NoError(err)
		suite.Require().Equal(2000.0, recvBt.Accounts["exchange"].Balances["ETH"])
		return nil
	}))
}
