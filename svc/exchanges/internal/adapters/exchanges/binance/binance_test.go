package binance

import (
	"context"
	"testing"
	"time"

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

	as.True(checkPairExistance(exch.Pairs, "ETH-USDC"))
	as.True(checkPairExistance(exch.Pairs, "FTM-USDC"))
	as.True(checkPairExistance(exch.Pairs, "BTC-USDC"))

	as.Equal(0.1, exch.Fees)

	as.WithinDuration(time.Now(), exch.LastSyncTime, time.Second)
}

func checkPairExistance(list []string, pair string) bool {
	for _, lp := range list {
		if pair == lp {
			return true
		}
	}

	return false
}
