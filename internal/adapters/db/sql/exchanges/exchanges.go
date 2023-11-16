package exchanges

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/internal/adapters/db/sql/exchanges/entities"
	"github.com/lerenn/cryptellation/internal/components/exchanges/ports/db"
	"github.com/lerenn/cryptellation/pkg/models/exchange"
	"gorm.io/gorm"
)

func (a *Adapter) CreateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	entities := make([]entities.Exchange, len(exchanges))
	for i, model := range exchanges {
		entities[i].FromModel(model)
	}

	err := a.db.Client.WithContext(ctx).Create(&entities).Error
	if err != nil {
		return fmt.Errorf("creating %+v: %w", exchanges, err)
	}

	return nil
}

func (a *Adapter) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	var ent []entities.Exchange
	if err := a.db.Client.WithContext(ctx).Preload("Pairs").Preload("Periods").Find(&ent, names).Error; err != nil {
		return nil, fmt.Errorf("reading %+v: %w", names, err)
	}

	models := make([]exchange.Exchange, len(ent))
	for i, entity := range ent {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func (a *Adapter) UpdateExchanges(ctx context.Context, exchanges ...exchange.Exchange) error {
	var entity entities.Exchange
	for _, model := range exchanges {
		entity.FromModel(model)

		if err := a.db.Client.WithContext(ctx).Updates(&entity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("updating %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("updating %+v: %w", exchanges, err)
		}

		if err := a.db.Client.WithContext(ctx).Model(&entity).Association("Pairs").Replace(entity.Pairs); err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("replacing pairs associations from %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("replacing pairs associations from %+v: %w", exchanges, err)
		}

		if err := a.db.Client.WithContext(ctx).Model(&entity).Association("Periods").Replace(entity.Periods); err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("replacing periods associations from %+v: %w", exchanges, db.ErrNotFound)
			}

			return fmt.Errorf("replacing periods associations from %+v: %w", exchanges, err)
		}
	}
	return nil
}

func (a *Adapter) DeleteExchanges(ctx context.Context, names ...string) error {
	for _, n := range names {
		if err := a.db.Client.WithContext(ctx).Delete(&entities.Exchange{
			Name: n,
		}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting %+v: %w", names, err)
		}

		err := a.db.Client.WithContext(ctx).
			Where("NOT EXISTS(SELECT NULL FROM exchanges_pairs ep WHERE ep.pair_symbol = symbol)").
			Delete(&entities.Pair{}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting unlinked pairs %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting unlinked pairs %+v: %w", names, err)
		}

		err = a.db.Client.WithContext(ctx).
			Where("NOT EXISTS(SELECT NULL FROM exchanges_periods ep WHERE ep.period_symbol = symbol)").
			Delete(&entities.Period{}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("deleting unlinked periods %+v: %w", names, db.ErrNotFound)
			}

			return fmt.Errorf("deleting unlinked periods %+v: %w", names, err)
		}
	}
	return nil
}
