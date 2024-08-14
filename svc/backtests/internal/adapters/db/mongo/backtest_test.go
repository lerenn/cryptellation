package backtests

import (
	"context"
	"errors"
	"testing"
	"time"

	"cryptellation/internal/config"
	"cryptellation/pkg/models/account"

	"cryptellation/svc/backtests/internal/app/ports/db"
	"cryptellation/svc/backtests/pkg/backtest"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	db.BacktestSuite
}

func (suite *BacktestSuite) SetupTest() {
	db, err := New(
		context.Background(),
		config.LoadMongo(&config.Mongo{
			Database: "cryptellation-backtests-integration-tests",
		}),
	)
	suite.Require().NoError(err)
	suite.DB = db
}

func (suite *BacktestSuite) TestPanicLock() {
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
	suite.Require().NoError(suite.DB.LockedBacktest(context.TODO(), bt.ID, func(bt *backtest.Backtest) error {
		panic(errors.New("PANIC"))
	}))
}
