package entities

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type SimpleMovingAverage struct {
	Exchange     string    `bson:"exchange"`
	Pair         string    `bson:"pair"`
	Period       string    `bson:"period"`
	PeriodNumber int       `bson:"period_number"`
	PriceType    string    `bson:"price_type"`
	Time         time.Time `bson:"time"`
	Price        float64
}

func (s *SimpleMovingAverage) FromModel(
	exchange, pair string,
	period period.Symbol,
	periodNb int,
	priceType candlestick.PriceType,
	t time.Time,
	price float64,
) {
	s.Exchange = exchange
	s.Pair = pair
	s.Period = period.String()
	s.PeriodNumber = periodNb
	s.PriceType = priceType.String()
	s.Time = t
	s.Price = price
}

func (s SimpleMovingAverage) ToModel() (exchange, pair, period string, periodNb int, priceType candlestick.PriceType, t time.Time, price float64) {
	return s.Exchange,
		s.Pair,
		s.Period,
		s.PeriodNumber,
		candlestick.PriceType(s.PriceType),
		s.Time,
		s.Price
}

func FromModelListToEntityList(
	exchange, pair string,
	period period.Symbol,
	periodNb int,
	priceType candlestick.PriceType,
	ts *timeserie.TimeSerie[float64],
) []SimpleMovingAverage {
	entities := make([]SimpleMovingAverage, 0, ts.Len())
	_ = ts.Loop(func(t time.Time, p float64) (bool, error) {
		sma := SimpleMovingAverage{}
		sma.FromModel(exchange, pair, period, periodNb, priceType, t, p)
		entities = append(entities, sma)
		return false, nil
	})

	return entities
}

func FromEntityListToModelList(entities []SimpleMovingAverage) (*timeserie.TimeSerie[float64], error) {
	ts := timeserie.New[float64]()
	for _, e := range entities {
		ts.Set(e.Time, e.Price)
	}

	return ts, nil
}
