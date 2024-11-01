package backtests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"

	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

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
	suite.Require().NoError(db.Reset(context.TODO()))
	suite.DB = db
}

func (suite *BacktestSuite) TestPanicLock() {
	bt := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Parameters{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceIsLow,
		},
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
