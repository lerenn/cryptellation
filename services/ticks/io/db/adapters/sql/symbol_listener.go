package sql

import (
	"context"
	"errors"

	"github.com/lerenn/cryptellation/services/ticks/io/db/adapters/sql/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (sqldb *DB) IncrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:    exchange,
		PairSymbol:  pairSymbol,
		Subscribers: 1, // Default if created
	}

	err := sqldb.client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return 1, sqldb.client.WithContext(ctx).Create(&sl).Error
	case err == nil:
		sl.Subscribers += 1
		return sl.Subscribers, sqldb.client.WithContext(ctx).Save(&sl).Error
	default:
		return 0, err
	}
}

func (sqldb *DB) DecrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := sqldb.client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	if err != nil {
		return 0, err
	}

	sl.Subscribers -= 1
	return sl.Subscribers, sqldb.client.WithContext(ctx).Save(&sl).Error
}

func (sqldb *DB) GetSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := sqldb.client.WithContext(ctx).Find(&sl).Error
	return sl.Subscribers, err
}

func (sqldb *DB) ClearSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) error {
	sl := entities.SymbolListener{
		Exchange:   exchange,
		PairSymbol: pairSymbol,
	}

	err := sqldb.client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sl).Error
	if err != nil {
		return err
	}

	sl.Subscribers = 0
	return sqldb.client.WithContext(ctx).Save(&sl).Error
}

func (sqldb *DB) ClearAllSymbolListenersCount(ctx context.Context) error {
	var sls []entities.SymbolListener

	err := sqldb.client.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sls).Error
	if err != nil {
		return err
	}

	for i := range sls {
		sls[i].Subscribers = 0
	}

	return sqldb.client.WithContext(ctx).Save(sls).Error
}
