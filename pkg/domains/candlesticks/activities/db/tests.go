package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"github.com/stretchr/testify/suite"
)

// CandlesticksSuite is the test suite for the candlesticks database.
type CandlesticksSuite struct {
	suite.Suite
	DB DB
}

// TestCreate tests the case where the candlesticks are created.
func (suite *CandlesticksSuite) TestCreate() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	t := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:       t,
		Open:       1,
		Low:        0.5,
		High:       2,
		Close:      1.5,
		Volume:     1000,
		Uncomplete: true,
	}))

	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t.Add(-time.Hour),
		End:      t.Add(time.Hour),
		Limit:    0,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(1, res.List.Data.Len())
	cs, _ := list.Data.Get(t)
	rcs, exists := res.List.Data.Get(t)
	suite.True(exists)
	suite.True(cs == rcs)
}

// TestCreateTwice tests the case where the candlesticks are created twice.
func (suite *CandlesticksSuite) TestCreateTwice() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	t := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t,
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)
	_, err = suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().Error(err)
}

// TestCreateWithNoTime tests the case where the candlestick to create does not
// have a time set.
func (suite *CandlesticksSuite) TestCreateWithNoTime() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))
	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().Error(err)
}

// TestRead tests the case where the candlesticks are read.
// TODO: Refactor this function
//
//nolint:funlen
func (suite *CandlesticksSuite) TestRead() {
	// Create targeted exchange, pair, period
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)

	t1 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t1,
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	t2 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:30:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t2,
		Open:   2,
		Low:    1,
		High:   4,
		Close:  3,
		Volume: 2000,
	}))

	t3 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:31:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t3,
		Open:   3,
		Low:    1.5,
		High:   6,
		Close:  4.5,
		Volume: 3000,
	}))

	t4 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:32:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t4,
		Open:   4,
		Low:    2,
		High:   8,
		Close:  6,
		Volume: 4000,
	}))

	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)

	// Create other exchange
	otherExchangeList := candlestick.NewList("exchange2", "ETH-USDC", period.M1)
	suite.Require().NoError(otherExchangeList.Set(candlestick.Candlestick{
		Time:   t2,
		Open:   1,
		Low:    1,
		High:   1,
		Close:  1,
		Volume: 1,
	}))
	_, err = suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: otherExchangeList,
	})
	suite.Require().NoError(err)

	// Create other pair
	otherPairList := candlestick.NewList("exchange", "BTC-USDC", period.M1)
	suite.Require().NoError(otherPairList.Set(candlestick.Candlestick{
		Time:   t2,
		Open:   2,
		Low:    2,
		High:   2,
		Close:  2,
		Volume: 2,
	}))
	_, err = suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: otherPairList,
	})
	suite.Require().NoError(err)

	// Create other period
	otherPeriodList := candlestick.NewList("exchange", "ETH-USDC", period.M15)
	suite.Require().NoError(otherPeriodList.Set(candlestick.Candlestick{
		Time:   t2,
		Open:   3,
		Low:    3,
		High:   3,
		Close:  3,
		Volume: 3,
	}))
	_, err = suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: otherPeriodList,
	})
	suite.Require().NoError(err)

	// Read only the two centered candles
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t2,
		End:      t3,
		Limit:    0,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(2, res.List.Data.Len())
	c, _ := list.Data.Get(t2)
	rc, exists := res.List.Data.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	c, _ = list.Data.Get(t3)
	rc, exists = res.List.Data.Get(t3)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	// Check others
	res, err = suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange2",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t2,
		End:      t3,
		Limit:    0,
	})
	suite.Require().NoError(err)
	c, _ = otherExchangeList.Data.Get(t2)
	rc, exists = res.List.Data.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	res, err = suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "BTC-USDC",
		Period:   period.M1,
		Start:    t2,
		End:      t3,
		Limit:    0,
	})
	suite.Require().NoError(err)
	c, _ = otherPairList.Data.Get(t2)
	rc, exists = res.List.Data.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	res, err = suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M15,
		Start:    t2,
		End:      t3,
		Limit:    0,
	})
	suite.Require().NoError(err)
	c, _ = otherPeriodList.Data.Get(t2)
	rc, exists = res.List.Data.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)
}

// TestReadLimit tests the case where the limit is set.
// TODO(lerenn): Refactor this function
//
//nolint:funlen
func (suite *CandlesticksSuite) TestReadLimit() {
	// Create targeted exchange, pair, period
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)

	t1 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t1,
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	t2 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:30:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t2,
		Open:   2,
		Low:    1,
		High:   4,
		Close:  3,
		Volume: 2000,
	}))

	t3 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:31:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t3,
		Open:   3,
		Low:    1.5,
		High:   6,
		Close:  4.5,
		Volume: 3000,
	}))

	t4 := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:32:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t4,
		Open:   4,
		Low:    2,
		High:   8,
		Close:  6,
		Volume: 4000,
	}))

	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)

	// Read only the 2 first
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t1,
		End:      t4,
		Limit:    2,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(2, res.List.Data.Len())
	c, _ := list.Data.Get(t1)
	rc, exists := res.List.Data.Get(t1)
	suite.True(exists)
	suite.Require().Equal(c, rc)

	c, _ = list.Data.Get(t2)
	rc, exists = res.List.Data.Get(t2)
	suite.True(exists)
	suite.Require().Equal(c, rc)
}

// TestReadEmpty tests the case where there is no candlestick to read.
func (suite *CandlesticksSuite) TestReadEmpty() {
	t, err := time.Parse(time.RFC3339, "1993-11-15T11:29:00Z")
	suite.Require().NoError(err)
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t.Add(-time.Hour),
		End:      t.Add(time.Hour),
		Limit:    0,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(0, res.List.Data.Len())
}

// TestUpdate tests the case where the candlestick to update exists.
func (suite *CandlesticksSuite) TestUpdate() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	t := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:       t,
		Open:       1,
		Low:        0.5,
		High:       2,
		Close:      1.5,
		Volume:     1000,
		Uncomplete: true,
	}))
	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)

	update := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(update.Set(candlestick.Candlestick{
		Time:       t,
		Open:       2,
		Low:        1,
		High:       4,
		Close:      3,
		Volume:     2000,
		Uncomplete: false,
	}))
	_, err = suite.DB.UpdateCandlesticksActivity(context.Background(), UpdateCandlesticksActivityParams{
		List: update,
	})
	suite.Require().NoError(err)
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    t.Add(-time.Hour),
		End:      t.Add(time.Hour),
		Limit:    0,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(1, res.List.Data.Len())
	cs, _ := update.Data.Get(t)
	rcs, exists := res.List.Data.Get(t)
	suite.True(exists)
	suite.Require().Equal(cs, rcs)
}

// TestUpdateInexistantTwice tests the case where the candlestick to update does
// not exist twice.
func (suite *CandlesticksSuite) TestUpdateInexistantTwice() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	t := utils.Must(time.Parse(time.RFC3339, "1993-11-15T11:29:00Z"))
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Time:   t,
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	_, err := suite.DB.UpdateCandlesticksActivity(context.Background(), UpdateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().Error(err)

	_, err = suite.DB.UpdateCandlesticksActivity(context.Background(), UpdateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().Error(err)
}

// TestUpdateWithNoTime tests the case where the candlestick to update does
// not have a time set.
func (suite *CandlesticksSuite) TestUpdateWithNoTime() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(list.Set(candlestick.Candlestick{
		Open:   1,
		Low:    0.5,
		High:   2,
		Close:  1.5,
		Volume: 1000,
	}))

	_, err := suite.DB.UpdateCandlesticksActivity(context.Background(), UpdateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().Error(err)
}

// TestDelete tests the case where the candlesticks are deleted.
func (suite *CandlesticksSuite) TestDelete() {
	list := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := 0; i < 10; i++ {
		suite.Require().NoError(list.Set(candlestick.Candlestick{
			Time:   time.Unix(int64(i*int(period.M1.Duration().Seconds())), 0),
			Open:   1 + float64(i),
			Low:    0.5 + float64(i),
			High:   2 + float64(i),
			Close:  1.5 + float64(i),
			Volume: 1000 * float64(i),
		}))
	}
	_, err := suite.DB.CreateCandlesticksActivity(context.Background(), CreateCandlesticksActivityParams{
		List: list,
	})
	suite.Require().NoError(err)

	// Remove half the data
	halfList := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(
		halfList.Set(candlestick.Candlestick{Time: time.Unix(int64(0*period.M1.Duration().Seconds()), 0)}))
	suite.Require().NoError(
		halfList.Set(candlestick.Candlestick{Time: time.Unix(int64(1*period.M1.Duration().Seconds()), 0)}))
	suite.Require().NoError(
		halfList.Set(candlestick.Candlestick{Time: time.Unix(int64(2*period.M1.Duration().Seconds()), 0)}))
	suite.Require().NoError(
		halfList.Set(candlestick.Candlestick{Time: time.Unix(int64(3*period.M1.Duration().Seconds()), 0)}))
	suite.Require().NoError(
		halfList.Set(candlestick.Candlestick{Time: time.Unix(int64(4*period.M1.Duration().Seconds()), 0)}))
	_, err = suite.DB.DeleteCandlesticksActivity(context.Background(), DeleteCandlesticksActivityParams{
		List: halfList,
	})
	suite.Require().NoError(err)

	// Check staying data
	tEnd := time.Unix(10*int64(period.M1.Duration().Seconds()), 0)
	tStart := time.Unix(5*int64(period.M1.Duration().Seconds()), 0)
	res, err := suite.DB.ReadCandlesticksActivity(context.Background(), ReadCandlesticksActivityParams{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    tStart,
		End:      tEnd,
		Limit:    0,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(list.Data.Len()-5, res.List.Data.Len())
	suite.Require().NoError(res.List.Loop(func(cs candlestick.Candlestick) (bool, error) {
		origCS, exists := list.Data.Get(cs.Time)
		suite.Require().True(exists)
		suite.Require().True(origCS.Equal(cs))
		return false, nil
	}))
}
