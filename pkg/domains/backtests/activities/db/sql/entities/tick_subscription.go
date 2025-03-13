package entities

import "github.com/lerenn/cryptellation/v1/pkg/models/tick"

// TickSubscription is the entity for a tick subscription.
type TickSubscription struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}

// ToModel converts the entity to a model.
func (ts TickSubscription) ToModel() tick.Subscription {
	return tick.Subscription{
		Exchange: ts.Exchange,
		Pair:     ts.Pair,
	}
}

// ToTickSubscriptionModels converts a slice of entities to a slice of models.
func ToTickSubscriptionModels(entities []TickSubscription) []tick.Subscription {
	models := make([]tick.Subscription, len(entities))
	for i, e := range entities {
		models[i] = e.ToModel()
	}
	return models
}

// FromTickSubscriptionModels converts a slice of models to a slice of entities.
func FromTickSubscriptionModels(models []tick.Subscription) []TickSubscription {
	entities := make([]TickSubscription, len(models))
	for i, m := range models {
		entities[i] = FromTickSubscriptionModel(m)
	}
	return entities
}

// FromTickSubscriptionModel converts a model to an entity.
func FromTickSubscriptionModel(m tick.Subscription) TickSubscription {
	return TickSubscription{
		Exchange: m.Exchange,
		Pair:     m.Pair,
	}
}
