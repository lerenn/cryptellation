package ticks

import (
	"context"
	"errors"

	"github.com/lerenn/cryptellation/internal/adapters/db/sql/ticks/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (a *Adapter) IncrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:    exchange,
		PairSymbol:  pairSymbol,
		Subscribers: 1, // Default if created
	}

	err := a.db.Client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return 1, a.db.Client.WithContext(ctx).Create(&sl).Error
	case err == nil:
		sl.Subscribers += 1
		return sl.Subscribers, a.db.Client.WithContext(ctx).Save(&sl).Error
	default:
		return 0, err
	}
}

func (a *Adapter) DecrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := a.db.Client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	if err != nil {
		return 0, err
	}

	sl.Subscribers -= 1
	return sl.Subscribers, a.db.Client.WithContext(ctx).Save(&sl).Error
}

func (a *Adapter) GetSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := a.db.Client.WithContext(ctx).Find(&sl).Error
	return sl.Subscribers, err
}

func (a *Adapter) ClearSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) error {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := a.db.Client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	if err != nil {
		return err
	}

	sl.Subscribers = 0
	return a.db.Client.WithContext(ctx).Save(&sl).Error
}

func (a *Adapter) ClearAllSymbolListenersCount(ctx context.Context) error {
	var sls []entities.SymbolListener

	err := a.db.Client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sls).Error
	if err != nil {
		return err
	}

	for i := range sls {
		sls[i].Subscribers = 0
	}

	return a.db.Client.WithContext(ctx).Save(sls).Error
}
