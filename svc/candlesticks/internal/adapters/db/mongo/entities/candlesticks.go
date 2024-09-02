package entities

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"golang.org/x/xerrors"
)

type Candlestick struct {
	Exchange   string    `bson:"exchange"`
	Pair       string    `bson:"pair"`
	Period     string    `bson:"period"`
	Time       time.Time `bson:"time"`
	Open       float64   `bson:"open"`
	High       float64   `bson:"high"`
	Low        float64   `bson:"low"`
	Close      float64   `bson:"close"`
	Volume     float64   `bson:"volume"`
	Uncomplete bool      `bson:"uncomplete"`
}

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
	c.Time = model.Time
	c.Open = model.Open
	c.High = model.High
	c.Low = model.Low
	c.Close = model.Close
	c.Volume = model.Volume
	c.Uncomplete = model.Uncomplete

	return nil
}

func (c Candlestick) ToModel() (exchange, pair, period string, model candlestick.Candlestick) {
	return c.Exchange,
		c.Pair,
		c.Period,
		candlestick.Candlestick{
			Time:       c.Time,
			Open:       c.Open,
			High:       c.High,
			Low:        c.Low,
			Close:      c.Close,
			Volume:     c.Volume,
			Uncomplete: c.Uncomplete,
		}
}

func FromModelListToEntityList(list *candlestick.List) ([]Candlestick, error) {
	entities := make([]Candlestick, 0, list.Len())
	err := list.Loop(func(cs candlestick.Candlestick) (bool, error) {
		c := Candlestick{}
		err := c.FromModel(list.Exchange, list.Pair, list.Period.String(), cs)
		entities = append(entities, c)
		return false, err
	})

	return entities, err
}

func FromEntityListToModelList(entities []Candlestick) (*candlestick.List, error) {
	if len(entities) == 0 {
		return nil, xerrors.New("no entities provided")
	}

	periodString := entities[0].Period
	per, err := period.FromString(periodString)
	if err != nil {
		return nil, fmt.Errorf("from candlestick entity to list: %w", err)
	}

	list := candlestick.NewList(entities[0].Exchange, entities[0].Pair, per)

	for _, e := range entities {
		if e.Exchange != list.Exchange {
			txt := fmt.Sprintf("incompatible exchanges for same list: %q and %q", e.Exchange, list.Exchange)
			return nil, xerrors.New(txt)
		}

		if e.Pair != list.Pair {
			txt := fmt.Sprintf("incompatible pair for same list: %q and %q", e.Pair, list.Pair)
			return nil, xerrors.New(txt)
		}

		if e.Period != list.Period.String() {
			txt := fmt.Sprintf("incompatible period for same list: %q and %q", e.Period, list.Period.String())
			return nil, xerrors.New(txt)
		}

		err := list.Set(candlestick.Candlestick{
			Time:       e.Time,
			Open:       e.Open,
			High:       e.High,
			Low:        e.Low,
			Close:      e.Close,
			Volume:     e.Volume,
			Uncomplete: e.Uncomplete,
		})
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
