package entities

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
)

type Exchange struct {
	Name         string
	Pairs        []string
	Periods      []string
	Fees         float64
	LastSyncTime time.Time
}

func ExchangeFromModel(model exchange.Exchange) Exchange {
	return Exchange{
		Name:         model.Name,
		Pairs:        model.Pairs,
		Periods:      model.Periods,
		Fees:         model.Fees,
		LastSyncTime: model.LastSyncTime,
	}
}

func (e Exchange) ToModel() exchange.Exchange {
	m := exchange.Exchange{
		Name:         e.Name,
		Pairs:        e.Pairs,
		Periods:      e.Periods,
		Fees:         e.Fees,
		LastSyncTime: e.LastSyncTime,
	}

	return m
}
