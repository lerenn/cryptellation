package db

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/stretchr/testify/suite"
)

type BacktestSuite struct {
	suite.Suite
	DB Port
}

func (suite *BacktestSuite) TestCreateRead() {
	bt := backtest.Backtest{
		ID:        uuid.New(),
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsLow,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), bt))
	rp, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Require().NoError(err, bt.ID.String())

	suite.Require().Equal(bt.ID, rp.ID)
	suite.Require().Len(rp.Accounts, 1)
	suite.Require().Len(rp.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(bt.Accounts["exchange"].Balances["DAI"], rp.Accounts["exchange"].Balances["DAI"])
}

func (suite *BacktestSuite) TestUpdate() {
	bt := backtest.Backtest{
		ID:        uuid.New(),
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsLow,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), bt))
	bt2 := backtest.Backtest{
		ID:        bt.ID,
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsClose,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	// Should be changes here
	suite.Require().NoError(suite.DB.UpdateBacktest(context.TODO(), bt2))
	rp, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Require().NoError(err)

	suite.Require().Equal(bt.ID, rp.ID)
	suite.Require().Equal(bt2.ID, rp.ID)
	suite.Require().Len(rp.Accounts, 1)
	suite.Require().Len(rp.Accounts["exchange2"].Balances, 1)
	suite.Require().Equal(bt2.Accounts["exchange2"].Balances["USDC"], rp.Accounts["exchange2"].Balances["USDC"])
}

func (suite *BacktestSuite) TestDelete() {
	bt := backtest.Backtest{
		ID:        uuid.New(),
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsLow,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), bt))
	suite.Require().NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
	_, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Error(err)
}

func (suite *BacktestSuite) TestDeleteInexistant() {
	bt := backtest.Backtest{
		ID:        uuid.New(),
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsLow,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), bt))
	suite.Require().NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
	suite.Require().NoError(suite.DB.DeleteBacktest(context.TODO(), bt))
}

func (suite *BacktestSuite) TestLock() {
	bt := backtest.Backtest{
		ID:        uuid.New(),
		StartTime: time.Unix(0, 0),
		CurrentCsTick: backtest.CurrentCsTick{
			Time:      time.Unix(60, 0),
			PriceType: candlestick.PriceTypeIsLow,
		},
		EndTime:             time.Unix(120, 0),
		PeriodBetweenEvents: period.M1,
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), bt))

	// Check that lock works even with panic
	// suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
	// 	panic(errors.New("PANIC"))
	// }))

	for i := 0; i < 10; i++ {
		suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
			return nil
		}), fmt.Sprintf("Lock/Unlock attempt #%d", i+1))
	}

	go func() {
		err := suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
			bt.Accounts["exchange"].Balances["ETH"] = 2000
			time.Sleep(time.Second)
			return nil
		})
		suite.Require().NoError(err)
	}()
	time.Sleep(100 * time.Millisecond)

	suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
		bt.Accounts["exchange"].Balances["ETH"] = 3000
		return nil
	}))

	rp, err := suite.DB.ReadBacktest(context.TODO(), bt.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(3000.0, rp.Accounts["exchange"].Balances["ETH"])
}
