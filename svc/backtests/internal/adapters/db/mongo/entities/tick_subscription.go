package entities

import "github.com/lerenn/cryptellation/pkg/models/event"

type TickSubscription struct {
	Exchange string `bson:"exchange"`
	Pair     string `bson:"pair"`
}

func (ts TickSubscription) ToModel() event.PricesSubscription {
	return event.PricesSubscription{
		Exchange: ts.Exchange,
		Pair:     ts.Pair,
	}
}

func ToTickSubscriptionModels(entities []TickSubscription) []event.PricesSubscription {
	models := make([]event.PricesSubscription, len(entities))
	for i, e := range entities {
		models[i] = e.ToModel()
	}
	return models
}

func FromTickSubscriptionModels(models []event.PricesSubscription) []TickSubscription {
	entities := make([]TickSubscription, len(models))
	for i, m := range models {
		entities[i] = FromTickSubscriptionModel(m)
	}
	return entities
}

func FromTickSubscriptionModel(m event.PricesSubscription) TickSubscription {
	return TickSubscription{
		Exchange: m.Exchange,
		Pair:     m.Pair,
	}
}
