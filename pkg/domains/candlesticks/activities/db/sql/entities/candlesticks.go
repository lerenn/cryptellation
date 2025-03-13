package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

// CandlestickData is the entity for a candlestick data.
type CandlestickData struct {
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     float64 `json:"volume"`
	Uncomplete bool    `json:"uncomplete"`
}

// Candlestick is the entity for a candlestick.
type Candlestick struct {
	Exchange string    `db:"exchange"`
	Pair     string    `db:"pair"`
	Period   string    `db:"period"`
	Time     time.Time `db:"time"`
	Data     []byte    `db:"data"`
}

// FromModel will convert a candlestick model to a candlestick entity.
func (c *Candlestick) FromModel(exchange, pair, period string, model candlestick.Candlestick) error {
	// Check that the time is not zero
	if model.Time.IsZero() {
		return fmt.Errorf("candlestick time is zero")
	}

	// List wise
	c.Exchange = exchange
	c.Pair = pair
	c.Period = period

	// Candlestick wise
	c.Time = model.Time.UTC()

	// Candlestick data
	dataByte, err := json.Marshal(&CandlestickData{
		Open:       model.Open,
		High:       model.High,
		Low:        model.Low,
		Close:      model.Close,
		Volume:     model.Volume,
		Uncomplete: model.Uncomplete,
	})
	if err != nil {
		return fmt.Errorf("from candlestick model to entity: %w", err)
	}
	c.Data = dataByte

	return nil
}

// ToModel will convert a candlestick entity to a candlestick model.
func (c Candlestick) ToModel() (exchange, pair, period string, model candlestick.Candlestick, err error) {
	// Candlestick data
	data := CandlestickData{}
	err = json.Unmarshal(c.Data, &data)
	if err != nil {
		err = fmt.Errorf("from candlestick entity to model: %w", err)
		return
	}

	// Return the model
	return c.Exchange,
		c.Pair,
		c.Period,
		candlestick.Candlestick{
			Time:       c.Time.UTC(),
			Open:       data.Open,
			High:       data.High,
			Low:        data.Low,
			Close:      data.Close,
			Volume:     data.Volume,
			Uncomplete: data.Uncomplete,
		}, nil
}

// FromEntitiesToMap will convert a list of candlestick entities to a map.
func FromEntitiesToMap(entities []Candlestick) []map[string]interface{} {
	maps := make([]map[string]interface{}, 0, len(entities))

	for _, e := range entities {
		maps = append(maps, map[string]interface{}{
			"exchange": e.Exchange,
			"pair":     e.Pair,
			"period":   e.Period,
			"time":     e.Time,
			"data":     e.Data,
		})
	}

	return maps
}

// FromModelListToEntityList will convert a candlestick model list to a candlestick entity list.
func FromModelListToEntityList(list *candlestick.List) ([]Candlestick, error) {
	entities := make([]Candlestick, 0, list.Data.Len())
	err := list.Loop(func(cs candlestick.Candlestick) (bool, error) {
		c := Candlestick{}
		err := c.FromModel(list.Metadata.Exchange, list.Metadata.Pair, list.Metadata.Period.String(), cs)
		entities = append(entities, c)
		return false, err
	})

	return entities, err
}

// FromEntityListToModelList will convert a candlestick entity list to a candlestick model list.
func FromEntityListToModelList(entities []Candlestick) (*candlestick.List, error) {
	if len(entities) == 0 {
		return nil, errors.New("no entities provided")
	}

	periodString := entities[0].Period
	per, err := period.FromString(periodString)
	if err != nil {
		return nil, fmt.Errorf("from candlestick entity to list: %w", err)
	}

	list := candlestick.NewList(entities[0].Exchange, entities[0].Pair, per)

	for _, e := range entities {
		if e.Exchange != list.Metadata.Exchange {
			txt := fmt.Sprintf("incompatible exchanges for same list: %q and %q", e.Exchange, list.Metadata.Exchange)
			return nil, errors.New(txt)
		}

		if e.Pair != list.Metadata.Pair {
			txt := fmt.Sprintf("incompatible pair for same list: %q and %q", e.Pair, list.Metadata.Pair)
			return nil, errors.New(txt)
		}

		if e.Period != list.Metadata.Period.String() {
			txt := fmt.Sprintf("incompatible period for same list: %q and %q", e.Period, list.Metadata.Period.String())
			return nil, errors.New(txt)
		}

		// Candlestick data
		data := CandlestickData{}
		err = json.Unmarshal(e.Data, &data)
		if err != nil {
			return nil, fmt.Errorf("from candlestick entity to model: %w", err)
		}

		err := list.Set(candlestick.Candlestick{
			Time:       e.Time.UTC(),
			Open:       data.Open,
			High:       data.High,
			Low:        data.Low,
			Close:      data.Close,
			Volume:     data.Volume,
			Uncomplete: data.Uncomplete,
		})
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
