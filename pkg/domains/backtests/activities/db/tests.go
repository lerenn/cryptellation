package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

// BacktestSuite is a suite of tests for the backtest activity database.
type BacktestSuite struct {
	suite.Suite
	DB DB
}

// TestCreateRead tests that creatin then reading a backtest activity works.
func (suite *BacktestSuite) TestCreateRead() {
	bt := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	resp, err := suite.DB.ReadBacktestActivity(context.Background(), ReadBacktestActivityParams{
		ID: bt.ID,
	})
	suite.Require().NoError(err, bt.ID.String())

	suite.Require().Equal(bt.ID, resp.Backtest.ID)
	suite.Require().Len(resp.Backtest.Accounts, 1)
	suite.Require().Len(resp.Backtest.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(bt.Accounts["exchange"].Balances["DAI"], resp.Backtest.Accounts["exchange"].Balances["DAI"])
	suite.Require().Equal(backtest.ModeIsFullOHLC, resp.Backtest.Parameters.Mode)
}

// TestList tests that listing backtests returns the correct number of backtests.
func (suite *BacktestSuite) TestList() {
	bt1 := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}
	bt2 := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"DAI": 1000,
				},
			},
		},
	}

	_, err := suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt1,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt2,
	})
	suite.Require().NoError(err)

	resp, err := suite.DB.ListBacktestsActivity(context.Background(), ListBacktestsActivityParams{})
	suite.Require().NoError(err)

	suite.Require().Len(resp.Backtests, 2)
	suite.Require().Equal(bt1.ID, resp.Backtests[0].ID)
	suite.Require().Equal(bt2.ID, resp.Backtests[1].ID)
}

// TestUpdate tests that updating a backtest works.
func (suite *BacktestSuite) TestUpdate() {
	bt := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	bt2 := backtest.Backtest{
		ID: bt.ID,
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsClose,
		},
		Accounts: map[string]account.Account{
			"exchange2": {
				Balances: map[string]float64{
					"USDC": 1500,
				},
			},
		},
	}
	// Should be changes here
	_, err = suite.DB.UpdateBacktestActivity(context.Background(), UpdateBacktestActivityParams{
		Backtest: bt2,
	})
	suite.Require().NoError(err)
	resp, err := suite.DB.ReadBacktestActivity(context.Background(), ReadBacktestActivityParams{
		ID: bt.ID,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(bt.ID, resp.Backtest.ID)
	suite.Require().Equal(bt2.ID, resp.Backtest.ID)
	suite.Require().Len(resp.Backtest.Accounts, 1)
	suite.Require().Len(resp.Backtest.Accounts["exchange2"].Balances, 1)
	suite.Require().Equal(bt2.Accounts["exchange2"].Balances["USDC"], resp.Backtest.Accounts["exchange2"].Balances["USDC"])
}

// TestDelete tests that deleting a backtest works.
func (suite *BacktestSuite) TestDelete() {
	bt := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.DeleteBacktestActivity(context.Background(), DeleteBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.ReadBacktestActivity(context.Background(), ReadBacktestActivityParams{
		ID: bt.ID,
	})
	suite.Error(err)
}

// TestDeleteInexistant tests that deleting an inexistant backtest does not return an error.
func (suite *BacktestSuite) TestDeleteInexistant() {
	bt := backtest.Backtest{
		ID: uuid.New(),
		Parameters: backtest.Settings{
			StartTime:   time.Unix(0, 0),
			EndTime:     time.Unix(120, 0),
			Mode:        backtest.ModeIsFullOHLC,
			PricePeriod: period.M1,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  time.Unix(60, 0),
			Price: candlestick.PriceTypeIsLow,
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"ETH": 1000,
				},
			},
		},
	}
	_, err := suite.DB.CreateBacktestActivity(context.Background(), CreateBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.DeleteBacktestActivity(context.Background(), DeleteBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.DeleteBacktestActivity(context.Background(), DeleteBacktestActivityParams{
		Backtest: bt,
	})
	suite.Require().NoError(err)
}
