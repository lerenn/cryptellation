package entities

import "github.com/lerenn/cryptellation/pkg/models/event"

type TickSubscription struct {
	ID           uint `gorm:"primaryKey"`
	BacktestID   uint
	ExchangeName string
	PairSymbol   string
}

func (ts TickSubscription) ToModel() event.TickSubscription {
	return event.TickSubscription{
		ID:           int(ts.ID),
		ExchangeName: ts.ExchangeName,
		PairSymbol:   ts.PairSymbol,
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
		ID:           uint(m.ID),
		BacktestID:   backtestID,
		ExchangeName: m.ExchangeName,
		PairSymbol:   m.PairSymbol,
	}
}
