package entities

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type SimpleMovingAverage struct {
	ExchangeName string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PairSymbol   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PeriodSymbol string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PeriodNumber int       `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	PriceType    string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Time         time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Price        float64
}

func (s *SimpleMovingAverage) FromModel(exchange, pair, period string, periodNb int, priceType candlestick.PriceType, t time.Time, price float64) {
	s.ExchangeName = exchange
	s.PairSymbol = pair
	s.PeriodSymbol = period
	s.PeriodNumber = periodNb
	s.PriceType = priceType.String()
	s.Time = t
	s.Price = price
}

func (s SimpleMovingAverage) ToModel() (exchange, pair, period string, periodNb int, priceType candlestick.PriceType, t time.Time, price float64) {
	return s.ExchangeName,
		s.PairSymbol,
		s.PeriodSymbol,
		s.PeriodNumber,
		candlestick.PriceType(s.PriceType),
		s.Time,
		s.Price
}

func FromModelListToEntityList(exchange, pair, period string, periodNb int, priceType candlestick.PriceType, ts *timeserie.TimeSerie[float64]) []SimpleMovingAverage {
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
