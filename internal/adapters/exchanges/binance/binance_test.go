package binance

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/internal/components/candlesticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

func TestBinanceSuite(t *testing.T) {
	suite.Run(t, new(BinanceSuite))
}

type BinanceSuite struct {
	suite.Suite
	service *Service
}

func (suite *BinanceSuite) SetupTest() {
	service, err := New()
	suite.Require().NoError(err)
	suite.service = service
}

func (suite *BinanceSuite) TestExchangeInfos() {
	as := suite.Require()

	exch, err := suite.service.Infos(context.TODO())
	suite.NoError(err)

	as.True(checkPairExistance(exch.PairsSymbols, "ETH-USDC"))
	as.True(checkPairExistance(exch.PairsSymbols, "FTM-USDC"))
	as.True(checkPairExistance(exch.PairsSymbols, "BTC-USDC"))

	as.Equal(0.1, exch.Fees)

	as.WithinDuration(time.Now(), exch.LastSyncTime, time.Second)
}

func checkPairExistance(list []string, pairSymbol string) bool {
	for _, lp := range list {
		if pairSymbol == lp {
			return true
		}
	}

	return false
}

func (suite *BinanceSuite) TestTicks() {
	tickChan, stopChan, err := suite.service.ListenSymbol("BTC-USDT")
	suite.Require().NoError(err)

	select {
	case recvTick := <-tickChan:
		suite.Require().Equal("BTC-USDT", recvTick.PairSymbol)
	case <-time.After(1 * time.Second):
		suite.Require().FailNow("Timeout")
	}

	stopChan <- struct{}{}
}

func (suite *BinanceSuite) TestGetCandlesticks() {
	p := "BTC-USDC"

	ts, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:00:00")
	suite.Require().NoError(err)

	te, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:05:00")
	suite.Require().NoError(err)

	cs, err := suite.service.GetCandlesticks(context.TODO(),
		exchanges.GetCandlesticksPayload{
			PairSymbol: p,
			Period:     period.M1,
			Limit:      2,
			Start:      ts,
			End:        te,
		})
	suite.Require().NoError(err)
	suite.Require().Equal(p, cs.PairSymbol)
	suite.Require().Equal(period.M1, cs.Period)

	expected := candlestick.Candlestick{
		Open:   16084.16,
		High:   16093.26,
		Low:    16084.16,
		Close:  16093.26,
		Volume: 0.344592,
	}

	suite.Require().Equal(2, cs.Len())
	rc, exists := cs.Get(ts)
	suite.Require().True(exists)
	suite.Require().Equal(expected, rc)
}

func (suite *BinanceSuite) TestGetCandlesticksWithZeroLimit() {
	p := "BTC-USDC"

	ts, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:00:00")
	suite.Require().NoError(err)

	te, err := time.Parse("2006/01/02 15:04:05", "2020/11/15 00:05:00")
	suite.Require().NoError(err)

	_, err = suite.service.GetCandlesticks(context.TODO(),
		exchanges.GetCandlesticksPayload{
			PairSymbol: p,
			Period:     period.M1,
			Limit:      0,
			Start:      ts,
			End:        te,
		})
	suite.Require().NoError(err)
}
