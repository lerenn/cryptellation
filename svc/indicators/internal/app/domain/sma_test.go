package domain

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/lerenn/cryptellation/svc/indicators/internal/app"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestSMASuite(t *testing.T) {
	suite.Run(t, new(SMASuite))
}

type SMASuite struct {
	suite.Suite
	app          app.Indicators
	db           *db.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *SMASuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.app = New(suite.db, suite.candlesticks)
}

func (suite *SMASuite) TestAllExistWithNoneInDB() {
	cs := candlestick.NewList("exchange", "ETC-USDT", period.M1)
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(0, 0), Close: 5})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(60, 0), Close: 7})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(120, 0), Close: 10})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(180, 0), Close: 10})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(240, 0), Close: 25})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(300, 0), Close: 25})
	cs.MustSet(candlestick.Candlestick{Time: time.Unix(360, 0), Close: 25})

	// Set expected calls
	suite.db.EXPECT().GetSMA(context.Background(), db.ReadSMAPayload{
		Exchange:     "exchange",
		Pair:         "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(180, 0),
		End:          time.Unix(360, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceIsClose,
	}).Return(timeserie.New[float64](), nil)

	suite.candlesticks.EXPECT().Read(context.Background(), candlesticks.ReadCandlesticksPayload{
		Exchange: "exchange",
		Pair:     "ETC-USDT",
		Period:   period.M1,
		Start:    utils.ToReference(time.Unix(0, 0)),
		End:      utils.ToReference(time.Unix(360, 0)),
	}).Return(cs, nil)

	suite.db.EXPECT().UpsertSMA(context.Background(), db.WriteSMAPayload{
		Exchange:     "exchange",
		Pair:         "ETC-USDT",
		Period:       period.M1,
		PeriodNumber: 3,
		PriceType:    candlestick.PriceIsClose,
		TimeSerie: timeserie.New[float64]().
			Set(time.Unix(180, 0), 9).
			Set(time.Unix(240, 0), 15).
			Set(time.Unix(300, 0), 20).
			Set(time.Unix(360, 0), 25),
	})

	// Run operation
	sma, err := suite.app.GetCachedSMA(context.Background(), app.GetCachedSMAPayload{
		Exchange:     "exchange",
		Pair:         "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(181, 0),
		End:          time.Unix(365, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceIsClose,
	})
	suite.Require().NoError(err)

	suite.Require().Equal(4, sma.Len())
	suite.Require().NoError(sma.Loop(func(t time.Time, v float64) (bool, error) {
		switch t.Unix() {
		case 180:
			suite.Require().Equal(9.0, v)
		case 240:
			suite.Require().Equal(15.0, v)
		case 300:
			suite.Require().Equal(20.0, v)
		case 360:
			suite.Require().Equal(25.0, v)
		default:
			suite.Fail("Unexpected time", t)
		}
		return false, nil
	}))
}
