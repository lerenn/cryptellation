package entities

import (
	"encoding/json"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/indicators/sma"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
)

// SimpleMovingAverageData is the entity for the simple moving average data.
type SimpleMovingAverageData struct {
	Price float64 `db:"price"`
}

// SimpleMovingAverage is the entity for the simple moving average.
type SimpleMovingAverage struct {
	Exchange     string    `db:"exchange"`
	Pair         string    `db:"pair"`
	Period       string    `db:"period"`
	PeriodNumber int       `db:"period_number"`
	PriceType    string    `db:"price_type"`
	Time         time.Time `db:"time"`
	Data         []byte    `db:"data"`
}

// FromModel converts the model to an entity.
func (s *SimpleMovingAverage) FromModel(p sma.Point) error {
	// Set the bytes
	dataByte, err := json.Marshal(SimpleMovingAverageData{
		Price: p.Price,
	})
	if err != nil {
		return err
	}

	// Set the values
	s.Exchange = p.Exchange
	s.Pair = p.Pair
	s.Period = p.Period.String()
	s.PeriodNumber = p.PeriodNb
	s.PriceType = p.PriceType.String()
	s.Time = p.Time.UTC()
	s.Data = dataByte

	return nil
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

	// Unmarshal the data
	data := SimpleMovingAverageData{}
	if err := json.Unmarshal(s.Data, &data); err != nil {
		return sma.Point{}, err
	}

	return sma.Point{
		Exchange:  s.Exchange,
		Pair:      s.Pair,
		Period:    per,
		PeriodNb:  s.PeriodNumber,
		PriceType: pt,
		Time:      s.Time.UTC(),
		Price:     data.Price,
	}, nil
}

// FromModelListToEntityList converts a timeserie to a list of entities.
func FromModelListToEntityList(
	exchange, pair string,
	period period.Symbol,
	periodNb int,
	priceType candlestick.PriceType,
	ts *timeserie.TimeSerie[float64],
) ([]SimpleMovingAverage, error) {
	entities := make([]SimpleMovingAverage, 0, ts.Len())
	err := ts.Loop(func(t time.Time, p float64) (bool, error) {
		point := SimpleMovingAverage{}
		if err := point.FromModel(sma.Point{
			Exchange:  exchange,
			Pair:      pair,
			Period:    period,
			PeriodNb:  periodNb,
			PriceType: priceType,
			Time:      t.UTC(),
			Price:     p,
		}); err != nil {
			return false, err
		}

		entities = append(entities, point)
		return false, nil
	})

	return entities, err
}

// FromEntityListToModelList converts a list of entities to a timeserie.
func FromEntityListToModelList(entities []SimpleMovingAverage) (*timeserie.TimeSerie[float64], error) {
	ts := timeserie.New[float64]()
	for _, e := range entities {
		// Unmarshal the data
		data := SimpleMovingAverageData{}
		if err := json.Unmarshal(e.Data, &data); err != nil {
			return nil, err
		}

		ts.Set(e.Time, data.Price)
	}

	return ts, nil
}

// FromEntitiesToMap converts a list of entities to a map.
func FromEntitiesToMap(entities []SimpleMovingAverage) []map[string]interface{} {
	maps := make([]map[string]interface{}, 0, len(entities))
	for _, e := range entities {
		maps = append(maps, map[string]interface{}{
			"exchange":      e.Exchange,
			"pair":          e.Pair,
			"period":        e.Period,
			"period_number": e.PeriodNumber,
			"price_type":    e.PriceType,
			"time":          e.Time.UTC(),
			"data":          e.Data,
		})
	}

	return maps
}
