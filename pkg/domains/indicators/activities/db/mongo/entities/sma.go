package entities

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/indicators/sma"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
)

// SimpleMovingAverage is the entity for the simple moving average.
type SimpleMovingAverage struct {
	Exchange     string    `bson:"exchange"`
	Pair         string    `bson:"pair"`
	Period       string    `bson:"period"`
	PeriodNumber int       `bson:"period_number"`
	PriceType    string    `bson:"price_type"`
	Time         time.Time `bson:"time"`
	Price        float64   `bson:"price"`
}

// FromModel converts the model to an entity.
func (s *SimpleMovingAverage) FromModel(p sma.Point) {
	s.Exchange = p.Exchange
	s.Pair = p.Pair
	s.Period = p.Period.String()
	s.PeriodNumber = p.PeriodNb
	s.PriceType = p.PriceType.String()
	s.Time = p.Time
	s.Price = p.Price
}

// ToModel converts the entity to a model.
func (s SimpleMovingAverage) ToModel() (sma.Point, error) {
	// Validate period
	per := period.Symbol(s.Period)
	if err := per.Validate(); err != nil {
		return sma.Point{}, err
	}

	// Validate price type
	pt := candlestick.PriceType(s.PriceType)
	if err := pt.Validate(); err != nil {
		return sma.Point{}, err
	}

	return sma.Point{
		Exchange:  s.Exchange,
		Pair:      s.Pair,
		Period:    per,
		PeriodNb:  s.PeriodNumber,
		PriceType: pt,
		Time:      s.Time,
		Price:     s.Price,
	}, nil
}

// FromModelListToEntityList converts a timeserie to a list of entities.
func FromModelListToEntityList(
	exchange, pair string,
	period period.Symbol,
	periodNb int,
	priceType candlestick.PriceType,
	ts *timeserie.TimeSerie[float64],
) []SimpleMovingAverage {
	entities := make([]SimpleMovingAverage, 0, ts.Len())
	_ = ts.Loop(func(t time.Time, p float64) (bool, error) {
		point := SimpleMovingAverage{}
		point.FromModel(sma.Point{
			Exchange:  exchange,
			Pair:      pair,
			Period:    period,
			PeriodNb:  periodNb,
			PriceType: priceType,
			Time:      t,
			Price:     p,
		})
		entities = append(entities, point)
		return false, nil
	})

	return entities
}

// FromEntityListToModelList converts a list of entities to a timeserie.
func FromEntityListToModelList(entities []SimpleMovingAverage) (*timeserie.TimeSerie[float64], error) {
	ts := timeserie.New[float64]()
	for _, e := range entities {
		ts.Set(e.Time, e.Price)
	}

	return ts, nil
}
