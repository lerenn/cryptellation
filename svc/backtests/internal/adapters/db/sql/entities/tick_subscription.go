package entities

import "github.com/lerenn/cryptellation/svc/backtests/pkg/event"

type TickSubscription struct {
	ID         uint `gorm:"primaryKey"`
	BacktestID uint
	Exchange   string
	Pair       string
}

func (ts TickSubscription) ToModel() event.TickSubscription {
	return event.TickSubscription{
		ID:       int(ts.ID),
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

func FromTickSubscriptionModels(backtestID uint, models []event.TickSubscription) []TickSubscription {
	entities := make([]TickSubscription, len(models))
	for i, m := range models {
		entities[i] = FromTickSubscriptionModel(backtestID, m)
	}
	return entities
}

func FromTickSubscriptionModel(backtestID uint, m event.TickSubscription) TickSubscription {
	return TickSubscription{
		ID:         uint(m.ID),
		BacktestID: backtestID,
		Exchange:   m.Exchange,
		Pair:       m.Pair,
	}
}
