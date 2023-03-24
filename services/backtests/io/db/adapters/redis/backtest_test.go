package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/pkg/account"
	"github.com/digital-feather/cryptellation/pkg/backtest"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/period"
	"github.com/digital-feather/cryptellation/services/backtests/io/db/test"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	test.BacktestSuite
}

func (suite *BacktestSuite) SetupTest() {
	db, err := New()
	suite.Require().NoError(err)
	suite.DB = db
}

func (suite *BacktestSuite) TestPanicLock() {
	bt := backtest.Backtest{
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
	suite.Require().NoError(suite.DB.CreateBacktest(context.TODO(), &bt))

	// Check that lock works even with panic
	suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
		panic(errors.New("PANIC"))
	}))
}
