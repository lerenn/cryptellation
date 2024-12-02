package binance

import (
	"testing"

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
	temporal temporalclient.Client
	service  *Activities
}

func (suite *BinanceSuite) SetupTest() {
	suite.temporal = temporal.NewMockClient(gomock.NewController(suite.T()))
	service, err := New(suite.temporal)
	suite.Require().NoError(err)
	suite.service = service
}

func (suite *BinanceSuite) TestTicks() {
	// TODO
}
