package backtest

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestBacktestSuite(t *testing.T) {
	suite.Run(t, new(BacktestSuite))
}

type BacktestSuite struct {
	suite.Suite
}

func (suite *BacktestSuite) TestMarshalUnMarshalBinary() {
	bt := Backtest{
		Parameters: Settings{
			StartTime: time.Unix(0, 0).UTC(),
			EndTime:   time.Unix(200, 0).UTC(),
		},
		CurrentCandlestick: CurrentCandlestick{
			Time: time.Unix(100, 0).UTC(),
		},
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"USDC": 10000,
				},
			},
		},
	}

	content, err := bt.MarshalBinary()
	suite.Require().NoError(err)
	nbt := Backtest{}
	suite.Require().NoError(nbt.UnmarshalBinary(content))
	suite.Require().Equal(bt, nbt)
}

func (suite *BacktestSuite) TestIncrementPriceID() {
	// TODO
}

func (suite *BacktestSuite) TestBacktestCreateWithModeFullOHLC() {
	params := Parameters{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"USDC": 10000,
				},
			},
		},
		StartTime:   time.Unix(0, 0).UTC(),
		EndTime:     nil,
		Mode:        utils.ToReference(ModeIsFullOHLC),
		PricePeriod: utils.ToReference(period.M1),
	}

	bt, err := New(params)
	suite.Require().NoError(err)
	suite.Require().Equal(ModeIsFullOHLC, bt.Parameters.Mode)
	suite.Require().Equal(candlestick.PriceTypeIsOpen, bt.CurrentCandlestick.Price)
}

func (suite *BacktestSuite) TestBacktestCreateWithModeCloseOHLC() {
	params := Parameters{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{
					"USDC": 10000,
				},
			},
		},
		StartTime:   time.Unix(0, 0).UTC(),
		EndTime:     nil,
		Mode:        utils.ToReference(ModeIsCloseOHLC),
		PricePeriod: utils.ToReference(period.M1),
	}

	bt, err := New(params)
	suite.Require().NoError(err)
	suite.Require().Equal(ModeIsCloseOHLC, bt.Parameters.Mode)
	suite.Require().Equal(candlestick.PriceTypeIsClose, bt.CurrentCandlestick.Price)
}

func (suite *BacktestSuite) TestBacktestSetNewTimeWithFullOHLCMode() {
	bt := Backtest{
		Parameters: Settings{
			StartTime: time.Unix(0, 0).UTC(),
			EndTime:   time.Unix(200, 0).UTC(),
			Mode:      ModeIsFullOHLC,
		},
		CurrentCandlestick: CurrentCandlestick{
			Time:  time.Unix(100, 0).UTC(),
			Price: candlestick.PriceTypeIsClose,
		},
	}

	bt.SetCurrentTime(time.Unix(150, 0).UTC())
	suite.Require().Equal(time.Unix(150, 0).UTC(), bt.CurrentCandlestick.Time)
	suite.Require().Equal(candlestick.PriceTypeIsOpen, bt.CurrentCandlestick.Price)
}

func (suite *BacktestSuite) TestBacktestSetNewTimeWithCloseOHLCMode() {
	bt := Backtest{
		Parameters: Settings{
			StartTime: time.Unix(0, 0).UTC(),
			EndTime:   time.Unix(200, 0).UTC(),
			Mode:      ModeIsCloseOHLC,
		},
		CurrentCandlestick: CurrentCandlestick{
			Time:  time.Unix(100, 0).UTC(),
			Price: candlestick.PriceTypeIsClose,
		},
	}

	bt.SetCurrentTime(time.Unix(150, 0).UTC())
	suite.Require().Equal(time.Unix(150, 0).UTC(), bt.CurrentCandlestick.Time)
	suite.Require().Equal(candlestick.PriceTypeIsClose, bt.CurrentCandlestick.Price)
}
