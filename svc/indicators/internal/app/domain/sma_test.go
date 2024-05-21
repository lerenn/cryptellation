package domain

import (
	"context"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
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
	candlesticks *client.MockClient
}

func (suite *SMASuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = client.NewMockClient(gomock.NewController(suite.T()))
	suite.app = New(suite.db, suite.candlesticks)
}

func (suite *SMASuite) TestAllExistWithNoneInDB() {
	cs := candlestick.NewList("exchange", "ETC-USDT", period.M1)
	cs.MustSet(time.Unix(0, 0), candlestick.Candlestick{Close: 5})
	cs.MustSet(time.Unix(60, 0), candlestick.Candlestick{Close: 7})
	cs.MustSet(time.Unix(120, 0), candlestick.Candlestick{Close: 10})
	cs.MustSet(time.Unix(180, 0), candlestick.Candlestick{Close: 10})
	cs.MustSet(time.Unix(240, 0), candlestick.Candlestick{Close: 25})
	cs.MustSet(time.Unix(300, 0), candlestick.Candlestick{Close: 25})
	cs.MustSet(time.Unix(360, 0), candlestick.Candlestick{Close: 25})

	// Set expected calls
	suite.db.EXPECT().GetSMA(context.Background(), db.ReadSMAPayload{
		Exchange:     "exchange",
		Pair:         "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(180, 0),
		End:          time.Unix(360, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	}).Return(timeserie.New[float64](), nil)

	suite.candlesticks.EXPECT().Read(context.Background(), client.ReadCandlesticksPayload{
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
		PriceType:    candlestick.PriceTypeIsClose,
		TimeSerie: timeserie.New[float64]().
			Set(time.Unix(180, 0), 9).
			Set(time.Unix(240, 0), 15).
			Set(time.Unix(300, 0), 20).
			Set(time.Unix(360, 0), 25),
	})

	// Run operation
	_, err := suite.app.GetCachedSMA(context.Background(), app.GetCachedSMAPayload{
		Exchange:     "exchange",
		Pair:         "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(180, 0),
		End:          time.Unix(360, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	})
	suite.Require().NoError(err)
}
