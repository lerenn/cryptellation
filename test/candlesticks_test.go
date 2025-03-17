package test

import (
	"context"
	"time"

	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

func (suite *EndToEndSuite) TestReadCandlesticks() {
	// WHEN requesting a candlesticks list

	res, err := suite.client.ListCandlesticks(context.Background(), api.ListCandlesticksWorkflowParams{
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
	suite.Require().Equal(2, res.List.Data.Len())
	_ = res.List.Loop(func(cs candlestick.Candlestick) (bool, error) {
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
