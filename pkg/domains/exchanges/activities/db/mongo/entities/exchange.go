package entities

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
)

// Exchange is the exchange entity.
type Exchange struct {
	Name         string
	Pairs        []string
	Periods      []string
	Fees         float64
	LastSyncTime time.Time
}

// ExchangeFromModel will convert the model to an entity.
func ExchangeFromModel(model exchange.Exchange) Exchange {
	return Exchange{
		Name:         model.Name,
		Pairs:        model.Pairs,
		Periods:      model.Periods,
		Fees:         model.Fees,
		LastSyncTime: model.LastSyncTime,
	}
}

// ToModel will convert the entity to a model.
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
