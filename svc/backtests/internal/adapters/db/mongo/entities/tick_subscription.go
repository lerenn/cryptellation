package entities

import "github.com/lerenn/cryptellation/pkg/models/event"

type TickSubscription struct {
	Exchange string `bson:"exchange"`
	Pair     string `bson:"pair"`
}

func (ts TickSubscription) ToModel() event.TickSubscription {
	return event.TickSubscription{
		Exchange: ts.Exchange,
		Pair:     ts.Pair,
	}
}

func ToTickSubscriptionModels(entities []TickSubscription) []event.TickSubscription {
	models := make([]event.TickSubscription, len(entities))
	for i, e := range entities {
		models[i] = e.ToModel()
	}
	return models
}

func FromTickSubscriptionModels(models []event.TickSubscription) []TickSubscription {
	entities := make([]TickSubscription, len(models))
	for i, m := range models {
		entities[i] = FromTickSubscriptionModel(m)
	}
	return entities
}

func FromTickSubscriptionModel(m event.TickSubscription) TickSubscription {
	return TickSubscription{
		Exchange: m.Exchange,
		Pair:     m.Pair,
	}
}
