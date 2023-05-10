package candlesticks

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lerenn/cryptellation/pkg/candlestick"
	"github.com/lerenn/cryptellation/pkg/period"
	db "github.com/lerenn/cryptellation/services/candlesticks/io/db"
	"github.com/lerenn/cryptellation/services/candlesticks/io/exchanges"
	"github.com/stretchr/testify/suite"
)

func TestGetCachedSuite(t *testing.T) {
	suite.Run(t, new(GetCachedSuite))
}

type GetCachedSuite struct {
	suite.Suite
	component Interface
	db        *db.MockPort
	exchange  *exchanges.MockPort
}

func (suite *GetCachedSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))
	suite.component = New(suite.db, suite.exchange)
}

func (suite *GetCachedSuite) setMocksForAllExistWithNoneInDB() context.Context {
	ctx := context.Background()

	// Set list that will be pulled from exchange and created in DB
	l := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1,
	})
	for i := int64(0); i < 100; i++ {
		suite.Require().NoError(l.Set(time.Unix(i*60, 0), candlestick.Candlestick{Open: float64(60 * i)}))
	}

	// Set first call to know how much candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(540, 0),
		uint(0),
	).Return(nil)

	// Set call to check which candlestick exists or not
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(5940, 0),
		uint(0),
	).Return(nil)

	// Set call for creating candlesticks in database
	suite.db.EXPECT().CreateCandlesticks(ctx, l).Return(nil)

	// Set call for getting the candlesticks from the database
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange:   "exchange",
			PairSymbol: "ETH-USDC",
			Period:     period.M1,
			Start:      time.Unix(0, 0),
			End:        time.Unix(5940, 0),
			Limit:      0,
		},
	).Return(l, nil)

	return ctx
}

func (suite *GetCachedSuite) TestAllExistWithNoneInDB() {
	ctx := suite.setMocksForAllExistWithNoneInDB()

	// When a request is made
	l, err := suite.component.GetCached(ctx, GetCachedPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        TimeOpt(time.Unix(0, 0)),
		End:          TimeOpt(time.Unix(540, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(10, l.Len())
	i := 0
	_ = l.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().Equal(float64(60*i), cs.Open)
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		i++

		return false, nil
	})
}

func (suite *GetCachedSuite) setMocksForNoneExistWithNoneInDB() context.Context {
	ctx := context.Background()

	l := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1,
	})

	// Set list that will be pulled from exchange and created in DB
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(540, 0),
		uint(0),
	).Return(nil)

	// Set call for getting the candlesticks from the database
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange:   "exchange",
			PairSymbol: "ETH-USDC",
			Period:     period.M1,
			Start:      time.Unix(0, 0),
			End:        time.Unix(5940, 0),
			Limit:      0,
		},
	).Return(l, nil)

	return ctx
}

func (suite *GetCachedSuite) TestNoneExistWithNoneInDB() {
	// Set list that will be pulled from exchange and created in DB
	ctx := suite.setMocksForNoneExistWithNoneInDB()

	// When a request is made
	l, err := suite.component.GetCached(ctx, GetCachedPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        TimeOpt(time.Unix(0, 0)),
		End:          TimeOpt(time.Unix(540, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(0, l.Len())
}

func (suite *GetCachedSuite) setMocksForFromDBAndService() context.Context {
	ctx := context.Background()

	id := candlestick.ListID{
		ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1,
	}

	dbl := candlestick.NewEmptyList(id)
	for i := int64(0); i < 10; i++ {
		suite.Require().NoError(dbl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 4321}))
	}

	exchl := candlestick.NewEmptyList(id)
	for i := int64(10); i < 110; i++ {
		suite.Require().NoError(exchl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}

	// Set list that will be pulled from exchange and created in DB
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(1140, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(*dbl, nil)
	})

	// Set first call to know how much candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(6540, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(*dbl, nil)
	})

	// Set call for creating candlesticks in database
	suite.db.EXPECT().CreateCandlesticks(ctx, exchl).Return(nil)

	// Set call for getting the candlesticks from the database
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange:   "exchange",
			PairSymbol: "ETH-USDC",
			Period:     period.M1,
			Start:      time.Unix(600, 0),
			End:        time.Unix(6540, 0),
			Limit:      0,
		},
	).Return(exchl, nil)

	return ctx
}

func (suite *GetCachedSuite) TestFromDBAndService() {
	ctx := suite.setMocksForFromDBAndService()

	// When a request is made
	el, err := suite.component.GetCached(ctx, GetCachedPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        TimeOpt(time.Unix(0, 0)),
		End:          TimeOpt(time.Unix(1140, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(20, el.Len())
	i := 0
	_ = el.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		if i < 10 {
			suite.Require().Equal(float64(4321), cs.Close, i)
		} else {
			suite.Require().Equal(float64(1234), cs.Close, i)
		}
		i++
		return false, nil
	})
}

func (suite *GetCachedSuite) setMocksForFromDBAndServiceWithUncomplete() context.Context {
	ctx := context.Background()

	id := candlestick.ListID{
		ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1,
	}

	dbl := candlestick.NewEmptyList(id)
	for i := int64(0); i < 10; i++ {
		suite.Require().NoError(dbl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 4321}))
	}
	suite.Require().NoError(dbl.Set(time.Unix(540, 0), candlestick.Candlestick{Close: 4321, Uncomplete: true}))

	exchl := candlestick.NewEmptyList(id)
	for i := int64(0); i < 100; i++ {
		suite.Require().NoError(exchl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}

	// Set list that will be pulled from exchange and created in DB
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(1140, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(*dbl, nil)
	})

	// Set first call to know how much candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewEmptyList(
			candlestick.ListID{ExchangeName: "exchange", PairSymbol: "ETH-USDC", Period: period.M1},
		),
		time.Unix(0, 0),
		time.Unix(5940, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(*dbl, nil)
	})

	// Set call for creating candlesticks in database
	createdl := candlestick.NewEmptyList(id)
	for i := int64(10); i < 100; i++ {
		suite.Require().NoError(createdl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.db.EXPECT().CreateCandlesticks(ctx, createdl).Return(nil)

	// Set call for getting the candlesticks from the database
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange:   "exchange",
			PairSymbol: "ETH-USDC",
			Period:     period.M1,
			Start:      time.Unix(0, 0),
			End:        time.Unix(5940, 0),
			Limit:      0,
		},
	).Return(exchl, nil)

	return ctx
}

func (suite *GetCachedSuite) TestFromDBAndServiceWithUncomplete() {
	ctx := suite.setMocksForFromDBAndServiceWithUncomplete()

	// When a request is made
	el, err := suite.component.GetCached(ctx, GetCachedPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        TimeOpt(time.Unix(0, 0)),
		End:          TimeOpt(time.Unix(1140, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(20, el.Len())
	i := 0
	_ = el.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		suite.Require().Equal(float64(1234), cs.Close, i)
		i++
		return false, nil
	})
}

func TimeOpt(t time.Time) *time.Time {
	return &t
}
