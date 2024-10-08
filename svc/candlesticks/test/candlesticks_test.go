package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"

	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

func (suite *EndToEndSuite) TestReadCandlesticks() {
	// WHEN requesting a candlesticks list
	list, err := suite.client.Read(context.Background(), client.ReadCandlesticksPayload{
		Exchange: "binance",
		Pair:     "ETH-USDT",
		Period:   period.H1,
		Start:    utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2022-01-01T00:00:00Z"))),
		End:      utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2022-01-01T03:00:00Z"))),
		Limit:    2,
	})

	// THEN the request is successful
	suite.Require().NoError(err)

	// AND the response contains the proper candlesticks

	t1 := utils.Must(time.Parse(time.RFC3339, "2022-01-01T00:00:00Z"))
	t2 := utils.Must(time.Parse(time.RFC3339, "2022-01-01T01:00:00Z"))
	suite.Require().Equal(2, list.Len())
	_ = list.Loop(func(cs candlestick.Candlestick) (bool, error) {
		switch {
		case cs.Time.Equal(t1):
			suite.Require().True(cs.Equal(candlestick.Candlestick{
				Time:   t1,
				Open:   3676.220000,
				High:   3730.000000,
				Low:    3676.220000,
				Close:  3723.040000,
				Volume: 9023.374,
			}))
		case cs.Time.Equal(t2):
			suite.Require().True(cs.Equal(candlestick.Candlestick{
				Time:   t2,
				Open:   3723.040000,
				High:   3748.450000,
				Low:    3714.100000,
				Close:  3724.890000,
				Volume: 8997.7569,
			}))
		default:
			suite.FailNow(cs.String()+"should not be there", cs.Time)
		}
		return false, nil
	})
}
