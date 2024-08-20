package domain

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"

	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app"
	db "github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetCachedSuite(t *testing.T) {
	suite.Run(t, new(GetCachedSuite))
}

type GetCachedSuite struct {
	suite.Suite
	component app.Candlesticks
	db        *db.MockPort
	exchange  *exchanges.MockPort
}

func (suite *GetCachedSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))
	suite.component = New(suite.db, suite.exchange)
}

func (suite *GetCachedSuite) TestAllExistWithNoneInDB() {
	ctx := context.Background()
	start, end := int64(0), int64(9)

	// Setting mocks -----------------------------------------------------------

	// Set first call to know how much candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(start*60, 0),
		time.Unix(end*60, 0),
		uint(0),
	).Return(nil)

	// Set call for getting the candlesticks from the database
	dStart, dEnd := getStartEndDownload(start, end)
	l := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := dStart; i <= dEnd; i++ {
		suite.Require().NoError(l.Set(time.Unix(i*60, 0), candlestick.Candlestick{Open: float64(60 * i)}))
	}
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange: "exchange",
			Pair:     "ETH-USDC",
			Period:   period.M1,
			Start:    time.Unix(dStart*60, 0),
			End:      time.Unix(dEnd*60, 0),
			Limit:    0,
		},
	).Return(l, nil)

	// Set call to check which candlestick exists or not
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(dStart*60, 0),
		time.Unix(dEnd*60, 0),
		uint(0),
	).Return(nil)

	// Set call for creating candlesticks in database
	suite.db.EXPECT().CreateCandlesticks(ctx, l).Return(nil)

	// Executing test ----------------------------------------------------------

	// When a request is made
	l, err := suite.component.GetCached(ctx, app.GetCachedPayload{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(start*60, 0)),
		End:      utils.ToReference(time.Unix(end*60, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(int(end-start)+1, l.Len())
	i := 0
	_ = l.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().Equal(float64(60*i), cs.Open)
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		i++

		return false, nil
	})
}

func (suite *GetCachedSuite) TestNoneExistWithNoneInDB() {
	ctx := context.Background()
	start, end := int64(0), int64(9)

	// Setting mocks -----------------------------------------------------------

	// Set list that will be pulled from exchange and created in DB
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(start*60, 0),
		time.Unix(end*60, 0),
		uint(0),
	).Return(nil)

	// Set call for getting the candlesticks from the database
	dStart, dEnd := getStartEndDownload(start, end)
	l := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	suite.Require().NoError(l.FillMissing(
		time.Unix(dStart*60, 0),
		time.Unix(dEnd*60, 0),
		candlestick.Candlestick{}))
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange: "exchange",
			Pair:     "ETH-USDC",
			Period:   period.M1,
			Start:    time.Unix(dStart*60, 0),
			End:      time.Unix(dEnd*60, 0),
			Limit:    0,
		},
	).Return(l, nil)

	// Set call to check which candlestick exists or not
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(dStart*60, 0),
		time.Unix(dEnd*60, 0),
		uint(0),
	).Return(nil)

	// Set call for creating candlesticks in database
	suite.db.EXPECT().CreateCandlesticks(ctx, l).Return(nil)

	// Executing test ----------------------------------------------------------

	// When a request is made
	l, err := suite.component.GetCached(ctx, app.GetCachedPayload{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(start*60, 0)),
		End:      utils.ToReference(time.Unix(end*60, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)

	// And that all candlesticks are empty
	suite.Require().Equal(int(end-start)+1, l.Len())
	suite.Require().NoError(l.Loop(func(_ time.Time, c candlestick.Candlestick) (bool, error) {
		suite.Require().Equal(c, candlestick.Candlestick{})
		return false, nil
	}))
}

func (suite *GetCachedSuite) TestFromDBAndService() {
	ctx := context.Background()
	start, end := int64(0), int64(74)

	// Set mocks ---------------------------------------------------------------

	// Set call to get candlesticks in database
	dbl := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := start; i < start+50; i++ {
		suite.Require().NoError(dbl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 4321}))
	}
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(start*0, 0),
		time.Unix(end*60, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(dbl, nil)
	})

	// Set call for getting missing candlesticks from exchange
	dStart, dEnd := getStartEndDownload(start+50, end)
	exchl := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := dStart; i <= dEnd; i++ {
		suite.Require().NoError(exchl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange: "exchange",
			Pair:     "ETH-USDC",
			Period:   period.M1,
			Start:    time.Unix(dStart*60, 0),
			End:      time.Unix(dEnd*60, 0),
			Limit:    0,
		},
	).Return(exchl, nil)

	// Set first call to know how much of the exchange candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(start*60, 0),
		time.Unix(dEnd*60, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(dbl, nil)
	})

	// Set call for creating candlesticks in database
	created := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := start + 50; i <= dEnd; i++ {
		suite.Require().NoError(created.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.db.EXPECT().CreateCandlesticks(ctx, created).Return(nil)

	// Set call for updating candlesticks in database
	updated := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := dStart; i < start+50; i++ {
		suite.Require().NoError(updated.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.db.EXPECT().UpdateCandlesticks(ctx, updated).Return(nil)

	// Execute the tests -------------------------------------------------------

	// When a request is made
	el, err := suite.component.GetCached(ctx, app.GetCachedPayload{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(start*60, 0)),
		End:      utils.ToReference(time.Unix(end*60, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(int(end-start)+1, el.Len())
	i := 0
	_ = el.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		if i < int(dStart) { // Present
			suite.Require().Equal(float64(4321), cs.Close, i)
		} else { // Updated + Non existant
			suite.Require().Equal(float64(1234), cs.Close, i)
		}
		i++
		return false, nil
	})
}

func (suite *GetCachedSuite) TestFromDBAndServiceWithUncomplete() {
	ctx := context.Background()
	start, end := int64(0), int64(19)
	uncompletePos := int64(9)

	// Set mocks ---------------------------------------------------------------

	// Set call to db
	dbl := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := start; i <= end; i++ {
		suite.Require().NoError(dbl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 4321}))
	}
	suite.Require().NoError(dbl.Set(time.Unix(uncompletePos*60, 0), candlestick.Candlestick{Close: 4321, Uncomplete: true}))
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(start*60, 0),
		time.Unix(end*60, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(dbl, nil)
	})

	// Set call for getting the candlesticks from the exchange
	dStart, dEnd := getStartEndDownload(uncompletePos, uncompletePos)
	exchl := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := dStart; i <= dEnd; i++ {
		suite.Require().NoError(exchl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.exchange.EXPECT().GetCandlesticks(
		ctx,
		exchanges.GetCandlesticksPayload{
			Exchange: "exchange",
			Pair:     "ETH-USDC",
			Period:   period.M1,
			Start:    time.Unix(dStart*60, 0),
			End:      time.Unix(dEnd*60, 0),
			Limit:    0,
		},
	).Return(exchl, nil)

	// Set first call to know how much candlestick there is in the database
	suite.db.EXPECT().ReadCandlesticks(
		ctx,
		candlestick.NewList("exchange", "ETH-USDC", period.M1),
		time.Unix(dStart*60, 0),
		time.Unix(dEnd*60, 0),
		uint(0),
	).DoAndReturn(func(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
		return cs.Merge(dbl, nil)
	})

	// Set call for creating candlesticks in database
	createdl := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := dStart; i <= dEnd; i++ {
		if i >= start && i <= end {
			continue
		}
		suite.Require().NoError(createdl.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.db.EXPECT().CreateCandlesticks(ctx, createdl).Return(nil)

	// Set call for updating candlesticks in database
	updated := candlestick.NewList("exchange", "ETH-USDC", period.M1)
	for i := start; i <= end; i++ {
		suite.Require().NoError(updated.Set(time.Unix(i*60, 0), candlestick.Candlestick{Close: 1234}))
	}
	suite.db.EXPECT().UpdateCandlesticks(ctx, updated).Return(nil)

	// Execute the tests -------------------------------------------------------

	// When a request is made
	el, err := suite.component.GetCached(ctx, app.GetCachedPayload{
		Exchange: "exchange",
		Pair:     "ETH-USDC",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(start*60, 0)),
		End:      utils.ToReference(time.Unix(end*60, 0)),
	})

	// Then everything went well
	suite.Require().NoError(err)
	suite.Require().Equal(int(end-start)+1, el.Len())
	i := 0
	_ = el.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		suite.Require().Equal(float64(1234), cs.Close, i)
		i++
		return false, nil
	})
}

func getStartEndDownload(start, end int64) (int64, int64) {
	count := end - start
	start = start - (MinimalRetrievedMissingCandlesticks-count)/2
	end = end + (MinimalRetrievedMissingCandlesticks-count)/2
	return int64(start), int64(end)
}
