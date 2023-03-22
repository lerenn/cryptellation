package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/internal/backtests/app/ports/db/test"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/pkg/types/period"
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
	bt := domain.Backtest{
		StartTime: time.Unix(0, 0),
		CurrentCsTick: domain.CurrentCsTick{
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
	suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *domain.Backtest) error {
		panic(errors.New("PANIC"))
	}))
}
