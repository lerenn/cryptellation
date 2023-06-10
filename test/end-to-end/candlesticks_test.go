package endToEnd

import (
	"context"
	"testing"
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/clients/go/nats"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestCandlesticksSuite(t *testing.T) {
	suite.Run(t, new(CandlesticksSuite))
}

type CandlesticksSuite struct {
	suite.Suite
	client client.Candlesticks
}

func (suite *CandlesticksSuite) SetupSuite() {
	// Get config
	cfg := config.LoadDefaultNATSConfig()
	cfg.OverrideFromEnv()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := nats.NewCandlesticks(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *CandlesticksSuite) TearDownSuite() {
	suite.client.Close()
}

func (suite *CandlesticksSuite) TestReadCandlesticks() {
	// WHEN requesting a candlesticks list
	list, err := suite.client.Read(context.Background(), client.ReadCandlesticksPayload{
		ExchangeName: "binance",
		PairSymbol:   "ETH-USDT",
		Period:       period.H1,
		Start:        utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2022-01-01T00:00:00Z"))),
		End:          utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2022-01-01T03:00:00Z"))),
		Limit:        2,
	})

	// THEN the request is successful
	suite.Require().NoError(err)

	// AND the response contains the proper candlesticks

	suite.Require().Equal(2, list.Len())
	_ = list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		switch {
		case t.Equal(utils.Must(time.Parse(time.RFC3339, "2022-01-01T00:00:00Z"))):
			suite.Require().True(cs.Equal(candlestick.Candlestick{
				Open:   3676.220000,
				High:   3730.000000,
				Low:    3676.220000,
				Close:  3723.040000,
				Volume: 9023.374,
			}))
		case t.Equal(utils.Must(time.Parse(time.RFC3339, "2022-01-01T01:00:00Z"))):
			suite.Require().True(cs.Equal(candlestick.Candlestick{
				Open:   3723.040000,
				High:   3748.450000,
				Low:    3714.100000,
				Close:  3724.890000,
				Volume: 8997.7569,
			}))
		default:
			suite.FailNow(cs.String()+"should not be there", t)
		}
		return false, nil
	})
}
