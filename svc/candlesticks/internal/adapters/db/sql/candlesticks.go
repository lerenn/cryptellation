package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/internal/adapters/db/sql/entities"
	"github.com/lerenn/cryptellation/svc/candlesticks/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

func (a *Adapter) CreateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := entities.FromModelListToEntityList(cs)
	return a.db.Client.WithContext(ctx).Create(&listCE).Error
}

func (a *Adapter) ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
	tx := a.db.Client.Where(`
		exchange = ? AND
		pair = ? AND
		period = ? AND
		time BETWEEN ? AND ?`,
		cs.Exchange,
		cs.Pair,
		cs.Period.String(),
		start, end)

	if limit != 0 {
		tx = tx.Limit(int(limit))
	}

	cse := []entities.Candlestick{}
	if err := tx.WithContext(ctx).Find(&cse).Error; err != nil {
		return err
	}

	for _, ce := range cse {
		_, _, _, t, m := ce.ToModel()
		if err := cs.Set(t, m); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) UpdateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := entities.FromModelListToEntityList(cs)
	for _, ce := range listCE {
		tx := a.db.Client.WithContext(ctx).
			Select("*").
			Model(&entities.Candlestick{}).
			Where("exchange = ?", ce.Exchange).
			Where("pair = ?", ce.Pair).
			Where("period = ?", ce.Period).
			Where("time = ?", ce.Time).
			Updates(ce)

		if tx.Error != nil {
			return fmt.Errorf("updating candlestick %q: %w", ce.Time, tx.Error)
		} else if tx.RowsAffected == 0 {
			return db.ErrNotFound
		}
	}

	return nil
}

func (a *Adapter) DeleteCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := entities.FromModelListToEntityList(cs)
	return a.db.Client.WithContext(ctx).Delete(&listCE).Error
}
