package backtests

import (
	"context"

	"github.com/lerenn/cryptellation/internal/adapters/db/sql/backtests/entities"
	"github.com/lerenn/cryptellation/internal/components/backtests/ports/db"
	"github.com/lerenn/cryptellation/pkg/models/backtest"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (a *Adapter) CreateBacktest(ctx context.Context, bt *backtest.Backtest) error {
	var count int64
	if err := a.db.Client.WithContext(ctx).Model(&entities.Backtest{}).Count(&count).Error; err != nil {
		return err
	}
	bt.ID = uint(count) + 1

	entity := entities.FromBacktestModel(*bt)
	return a.db.Client.WithContext(ctx).Create(&entity).Error
}

func (a *Adapter) readBacktest(ctx context.Context, id uint, clauses []clause.Expression) (backtest.Backtest, error) {
	e := entities.Backtest{
		ID: id,
	}

	tx := a.db.Client.
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

func (a *Adapter) ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error) {
	return a.readBacktest(ctx, id, nil)
}

func (a *Adapter) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	entity := entities.FromBacktestModel(bt)

	return a.newTransaction(func() error {
		// Remove all balances before updating
		tx := a.db.Client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.Balance{})
		if tx.Error != nil {
			return tx.Error
		}

		// Remove all orders before updating
		tx = a.db.Client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.Order{})
		if tx.Error != nil {
			return tx.Error
		}

		// Remove all tick subscriptions before updating
		tx = a.db.Client.WithContext(ctx).Where("backtest_id = ?", bt.ID).Delete(&entities.TickSubscription{})
		if tx.Error != nil {
			return tx.Error
		}

		return a.db.Client.WithContext(ctx).Save(&entity).Error
	})
}

func (a *Adapter) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	entity := entities.FromBacktestModel(bt)
	return a.db.Client.WithContext(ctx).Delete(&entity).Error

}

func (a *Adapter) newTransaction(fn func() error) error {
	return a.db.Client.Transaction(func(tx *gorm.DB) error {
		// Set client as transaction and defer removal
		originalClient := a.db.Client
		a.db.Client = tx
		defer func() {
			a.db.Client = originalClient
		}()

		return fn()
	})
}

func (a *Adapter) LockedBacktest(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
	return a.newTransaction(func() error {
		bt, err := a.readBacktest(ctx, id, []clause.Expression{
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

		return a.UpdateBacktest(ctx, bt)
	})
}
