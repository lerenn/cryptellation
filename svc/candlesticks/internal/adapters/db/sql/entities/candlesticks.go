package entities

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"golang.org/x/xerrors"
)

type Candlestick struct {
	Exchange   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Pair       string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Period     string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Time       time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Uncomplete bool
}

func (c *Candlestick) FromModel(exchange, pair, period string, t time.Time, model candlestick.Candlestick) {
	c.Exchange = exchange
	c.Pair = pair
	c.Period = period
	c.Time = t

	c.Open = model.Open
	c.High = model.High
	c.Low = model.Low
	c.Close = model.Close
	c.Volume = model.Volume
	c.Uncomplete = model.Uncomplete
}

func (c Candlestick) ToModel() (exchange, pair, period string, t time.Time, model candlestick.Candlestick) {
	return c.Exchange,
		c.Pair,
		c.Period,
		c.Time,
		candlestick.Candlestick{
			Open:       c.Open,
			High:       c.High,
			Low:        c.Low,
			Close:      c.Close,
			Volume:     c.Volume,
			Uncomplete: c.Uncomplete,
		}
}

func FromModelListToEntityList(list *candlestick.List) []Candlestick {
	entities := make([]Candlestick, 0, list.Len())
	_ = list.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		c := Candlestick{}
		c.FromModel(list.Exchange, list.Pair, list.Period.String(), t, cs)
		entities = append(entities, c)
		return false, nil
	})

	return entities
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

		err := list.Set(e.Time, candlestick.Candlestick{
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
