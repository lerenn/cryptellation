package binance

import (
	"testing"

	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/temporal"
	"github.com/stretchr/testify/suite"
	temporalclient "go.temporal.io/sdk/client"
	"go.uber.org/mock/gomock"
)

func TestBinanceSuite(t *testing.T) {
	suite.Run(t, new(BinanceSuite))
}

type BinanceSuite struct {
	suite.Suite
	temporal   temporalclient.Client
	activities exchanges.Exchanges
}

func (suite *BinanceSuite) SetupTest() {
	suite.temporal = temporal.NewMockClient(gomock.NewController(suite.T()))
	acts, err := New(suite.temporal)
	suite.Require().NoError(err)
	suite.activities = acts
}

func (suite *BinanceSuite) TestTicks() {
	// TODO(#63): Implement this test
}
