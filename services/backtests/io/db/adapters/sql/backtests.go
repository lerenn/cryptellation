package sql

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/io/db"
	"github.com/digital-feather/cryptellation/services/backtests/io/db/adapters/sql/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (sqldb *DB) CreateBacktest(ctx context.Context, bt *backtest.Backtest) error {
	var count int64
	if err := sqldb.client.WithContext(ctx).Model(&entities.Backtest{}).Count(&count).Error; err != nil {
		return err
	}
	bt.ID = uint(count) + 1

	entity := entities.FromBacktestModel(*bt)
	return sqldb.client.WithContext(ctx).Create(&entity).Error
}

func (sqldb *DB) readBacktest(ctx context.Context, id uint, clauses []clause.Expression) (backtest.Backtest, error) {
	e := entities.Backtest{
		ID: id,
	}

	tx := sqldb.client.
		WithContext(ctx).
		Clauses(clauses...).
		Preload("Balances").
		Preload("Orders").
		Preload("TickSubscriptions").
		Find(&e)
	if tx.Error != nil {
		return backtest.Backtest{}, tx.Error
	} else if tx.RowsAffected == 0 {
		return backtest.Backtest{}, db.ErrRecordNotFound
	}

	return e.ToModel()
}

func (sqldb *DB) ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error) {
	return sqldb.readBacktest(ctx, id, nil)
}

func (sqldb *DB) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	entity := entities.FromBacktestModel(bt)

	return sqldb.newTransaction(func() error {
		// Remove all balances before updating
		tx := sqldb.client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.Balance{})
		if tx.Error != nil {
			return tx.Error
		}

		// Remove all orders before updating
		tx = sqldb.client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.Order{})
		if tx.Error != nil {
			return tx.Error
		}

		// Remove all tick subscriptions before updating
		tx = sqldb.client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.TickSubscription{})
		if tx.Error != nil {
			return tx.Error
		}

		return sqldb.client.WithContext(ctx).Save(&entity).Error
	})
}

func (sqldb *DB) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	entity := entities.FromBacktestModel(bt)
	return sqldb.client.WithContext(ctx).Delete(&entity).Error

}

func (sqldb *DB) newTransaction(fn func() error) error {
	return sqldb.client.Transaction(func(tx *gorm.DB) error {
		// Set client as transaction and defer removal
		originalClient := sqldb.client
		sqldb.client = tx
		defer func() {
			sqldb.client = originalClient
		}()

		return fn()
	})
}

func (sqldb *DB) LockedBacktest(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
	return sqldb.newTransaction(func() error {
		bt, err := sqldb.readBacktest(ctx, id, []clause.Expression{
			clause.Locking{
				Strength: "UPDATE",
			},
		})
		if err != nil {
			return err
		}

		if err = fn(&bt); err != nil {
			return err
		}

		return sqldb.UpdateBacktest(ctx, bt)
	})
}
