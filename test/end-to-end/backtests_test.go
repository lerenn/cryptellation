package endToEnd

import (
	"testing"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/clients/go/nats"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestBacktestsSuite(t *testing.T) {
	suite.Run(t, new(BacktestsSuite))
}

type BacktestsSuite struct {
	suite.Suite
	client client.Backtests
}

func (suite *BacktestsSuite) SetupSuite() {
	// Get config
	cfg := config.LoadDefaultNATSConfig()
	cfg.OverrideFromEnv()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := nats.NewBacktests(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *BacktestsSuite) TearDownSuite() {
	suite.client.Close()
}

func (suite *BacktestsSuite) TestCreationDeletion() {
	// suite.client.
}
