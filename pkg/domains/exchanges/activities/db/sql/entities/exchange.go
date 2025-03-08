package entities

import (
	"encoding/json"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
)

// ExchangeData is the exchange data entity.
type ExchangeData struct {
	Pairs        []string  `json:"pairs"`
	Periods      []string  `json:"periods"`
	Fees         float64   `json:"fees"`
	LastSyncTime time.Time `json:"last_sync_time"`
}

// Exchange is the exchange entity.
type Exchange struct {
	Name string `db:"name"`
	Data []byte `db:"data"`
}

// ExchangeFromModel will convert the model to an entity.
func ExchangeFromModel(model exchange.Exchange) (Exchange, error) {
	dataByte, err := json.Marshal(&ExchangeData{
		Pairs:        model.Pairs,
		Periods:      model.Periods,
		Fees:         model.Fees,
		LastSyncTime: model.LastSyncTime,
	})
	if err != nil {
		return Exchange{}, err
	}

	return Exchange{
		Name: model.Name,
		Data: dataByte,
	}, nil
}

// ToModel will convert the entity to a model.
func (e Exchange) ToModel() exchange.Exchange {
	data := ExchangeData{}
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return exchange.Exchange{}
	}

	m := exchange.Exchange{
		Name:         e.Name,
		Pairs:        data.Pairs,
		Periods:      data.Periods,
		Fees:         data.Fees,
		LastSyncTime: data.LastSyncTime,
	}

	return m
}

// FromEntitiesToMap will convert the entities to a map.
func FromEntitiesToMap(exchanges []Exchange) []map[string]interface{} {
	m := make([]map[string]interface{}, len(exchanges))
	for i, e := range exchanges {
		m[i] = map[string]interface{}{
			"name": e.Name,
			"data": e.Data,
		}
	}

	return m
}
