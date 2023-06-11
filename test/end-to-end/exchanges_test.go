package endToEnd

import (
	"context"
	"testing"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	suite.Suite
	client client.Exchanges
}

func (suite *ExchangesSuite) SetupSuite() {
	// Get config
	cfg := config.LoadDefaultNATSConfig()
	cfg.OverrideFromEnv()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := nats.NewExchanges(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *ExchangesSuite) TearDownSuite() {
	suite.client.Close()
}

func (suite *ExchangesSuite) TestReadExchanges() {
	// WHEN requesting a exchanges list
	list, err := suite.client.Read(context.Background(), "binance")

	// THEN the request is successful
	suite.Require().NoError(err)

	// AND the response contains the proper exchanges
	suite.Require().Len(list, 1)
	suite.Require().Equal("binance", list[0].Name)

	l := []string{"D1", "D3", "H1", "H12", "H2", "H4", "H6", "H8", "M1", "M15", "M3", "M30", "M5", "W1"}
	for i, s := range l {
		suite.Require().Contains(list[0].PeriodsSymbols, s, i)
	}
}
