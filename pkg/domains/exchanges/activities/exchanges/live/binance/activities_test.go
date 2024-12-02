package binance

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/domains/exchanges/activities/exchanges"
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
	service, err := New()
	suite.Require().NoError(err)
	suite.service = service
}

func (suite *BinanceSuite) TestExchangeInfos() {
	r, err := suite.service.GetExchangeActivity(context.TODO(), exchanges.GetExchangeActivityParams{})
	suite.NoError(err)

	suite.Require().True(checkPairExistance(r.Exchange.Pairs, "ETH-USDC"))
	suite.Require().True(checkPairExistance(r.Exchange.Pairs, "FTM-USDC"))
	suite.Require().True(checkPairExistance(r.Exchange.Pairs, "BTC-USDC"))

	suite.Require().Equal(0.1, r.Exchange.Fees)

	suite.Require().WithinDuration(time.Now(), r.Exchange.LastSyncTime, time.Second)
}

func checkPairExistance(list []string, pair string) bool {
	for _, lp := range list {
		if pair == lp {
			return true
		}
	}

	return false
}

func (suite *BinanceSuite) TestExchangeNames() {
	r, err := suite.service.ListExchangesActivity(context.TODO(), exchanges.ListExchangesActivityParams{})
	suite.NoError(err)

	suite.Require().Contains(r.List, "Binance")
}
