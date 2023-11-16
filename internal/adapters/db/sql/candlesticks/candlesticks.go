package candlesticks

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/internal/adapters/db/sql/candlesticks/entities"
	"github.com/lerenn/cryptellation/internal/components/candlesticks/ports/db"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

func (a *Adapter) CreateCandlesticks(ctx context.Context, cs *candlestick.List) error {
	listCE := entities.FromModelListToEntityList(cs)
	return a.db.Client.WithContext(ctx).Create(&listCE).Error
}

func (a *Adapter) ReadCandlesticks(ctx context.Context, cs *candlestick.List, start, end time.Time, limit uint) error {
	tx := a.db.Client.Where(`
		exchange_name = ? AND
		pair_symbol = ? AND
		period_symbol = ? AND
		time BETWEEN ? AND ?`,
		cs.ExchangeName,
		cs.PairSymbol,
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
			Where("exchange_name = ?", ce.ExchangeName).
			Where("pair_symbol = ?", ce.PairSymbol).
			Where("period_symbol = ?", ce.PeriodSymbol).
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
