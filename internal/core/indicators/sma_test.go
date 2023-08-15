package indicators

import (
	"context"
	"testing"
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/clients/go/mock"
	"github.com/lerenn/cryptellation/internal/core/indicators/ports/db"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestSMASuite(t *testing.T) {
	suite.Run(t, new(SMASuite))
}

type SMASuite struct {
	suite.Suite
	app          Interface
	db           *db.MockPort
	candlesticks *mock.MockCandlesticks
}

func (suite *SMASuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = mock.NewMockCandlesticks(gomock.NewController(suite.T()))
	suite.app = New(suite.db, suite.candlesticks)
}

func (suite *SMASuite) TestAllExistWithNoneInDB() {
	cs := candlestick.NewEmptyList("exchange", "ETC-USDT", period.M1)
	cs.MustSet(time.Unix(0, 0), candlestick.Candlestick{Close: 5})
	cs.MustSet(time.Unix(60, 0), candlestick.Candlestick{Close: 7})
	cs.MustSet(time.Unix(120, 0), candlestick.Candlestick{Close: 10})
	cs.MustSet(time.Unix(180, 0), candlestick.Candlestick{Close: 10})
	cs.MustSet(time.Unix(240, 0), candlestick.Candlestick{Close: 25})
	cs.MustSet(time.Unix(300, 0), candlestick.Candlestick{Close: 25})
	cs.MustSet(time.Unix(360, 0), candlestick.Candlestick{Close: 25})

	// Set expected calls
	suite.db.EXPECT().GetSMA(context.Background(), db.ReadSMAPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(180, 0),
		End:          time.Unix(360, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	}).Return(timeserie.New[float64](), nil)

	suite.candlesticks.EXPECT().Read(context.Background(), client.ReadCandlesticksPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETC-USDT",
		Period:       period.M1,
		Start:        utils.ToReference(time.Unix(0, 0)),
		End:          utils.ToReference(time.Unix(360, 0)),
	}).Return(cs, nil)

	suite.db.EXPECT().UpsertSMA(context.Background(), db.WriteSMAPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETC-USDT",
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
	_, err := suite.app.GetCachedSMA(context.Background(), GetCachedSMAPayload{
		ExchangeName: "exchange",
		PairSymbol:   "ETC-USDT",
		Period:       period.M1,
		Start:        time.Unix(180, 0),
		End:          time.Unix(360, 0),
		PeriodNumber: 3,
		PriceType:    candlestick.PriceTypeIsClose,
	})
	suite.Require().NoError(err)
}
