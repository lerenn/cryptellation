package binance

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestBinanceSuite(t *testing.T) {
	suite.Run(t, new(BinanceSuite))
}

type BinanceSuite struct {
	suite.Suite
	service *Activities
}

func (suite *BinanceSuite) SetupTest() {
	service, err := New(config.LoadBinanceTest())
	suite.Require().NoError(err)
	suite.service = service
}

func (suite *BinanceSuite) TestGetCandlesticks() {
	p := "BTC-USDC"

	ts := utils.Must(time.Parse("2006/01/02 15:04:05", "2020/11/15 00:00:00"))
	te := utils.Must(time.Parse("2006/01/02 15:04:05", "2020/11/15 00:05:00"))

	res, err := suite.service.GetCandlesticksActivity(context.TODO(),
		exchanges.GetCandlesticksActivityParams{
			Pair:   p,
			Period: period.M1,
			Limit:  2,
			Start:  ts,
			End:    te,
		})
	suite.Require().NoError(err)
	suite.Require().Equal(p, res.List.Metadata.Pair)
	suite.Require().Equal(period.M1, res.List.Metadata.Period)

	expected := candlestick.Candlestick{
		Time:   ts,
		Open:   16084.16,
		High:   16093.26,
		Low:    16084.16,
		Close:  16093.26,
		Volume: 0.344592,
	}

	suite.Require().Equal(2, res.List.Data.Len())
	rc, exists := res.List.Data.Get(ts)
	suite.Require().True(exists)
	suite.Require().True(expected.Equal(rc))
}

func (suite *BinanceSuite) TestGetCandlesticksWithZeroLimit() {
	p := "BTC-USDC"

	ts, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:00:00")
	suite.Require().NoError(err)

	te, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:05:00")
	suite.Require().NoError(err)

	_, err = suite.service.GetCandlesticksActivity(context.TODO(),
		exchanges.GetCandlesticksActivityParams{
			Pair:   p,
			Period: period.M1,
			Limit:  0,
			Start:  ts,
			End:    te,
		})
	suite.Require().NoError(err)
}
