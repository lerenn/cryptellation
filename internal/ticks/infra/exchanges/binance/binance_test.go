package binance

import (
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
