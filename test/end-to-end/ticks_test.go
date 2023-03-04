package endToEnd

import (
	"context"
	"testing"
	"time"

	ticks "github.com/digital-feather/cryptellation/internal/ticks/ctrl/nats"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestTicksSuite(t *testing.T) {
	suite.Run(t, new(TicksSuite))
}

type TicksSuite struct {
	suite.Suite
	client ticks.Client
}

func (suite *TicksSuite) SetupSuite() {
	// Get config
	cfg := config.LoadDefaultNATSConfig()
	cfg.OverrideFromEnv()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := ticks.New(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *TicksSuite) TearDownSuite() {
	suite.client.Close()
}

func (suite *TicksSuite) TestListen() {
	// Register listener
	err := suite.client.Register(context.Background(), ticks.TicksFilterPayload{
		ExchangeName: "binance",
		PairSymbol:   "BTC-USDT",
	})
	suite.Require().NoError(err)

	// Listen to ticks
	ch, err := suite.client.Listen(context.Background(), ticks.TicksFilterPayload{
		ExchangeName: "binance",
		PairSymbol:   "BTC-USDT",
	})
	suite.Require().NoError(err)

	// Check that ticks are correct
	for i := 0; i < 3; i++ {
		t := <-ch
		suite.Require().Equal("binance", t.Exchange)
		suite.Require().Equal("BTC-USDT", t.PairSymbol)
		suite.Require().NotEqual(0, t.Price)
		suite.Require().WithinDuration(time.Now(), t.Time, time.Second)
	}

	// Unregister listener
	err = suite.client.Unregister(context.Background(), ticks.TicksFilterPayload{
		ExchangeName: "binance",
		PairSymbol:   "BTC-USDT",
	})
	suite.Require().NoError(err)
}
