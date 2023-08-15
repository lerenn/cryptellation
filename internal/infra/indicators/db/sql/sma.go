package sql

import (
	"context"

	"github.com/lerenn/cryptellation/internal/core/indicators/ports/db"
	"github.com/lerenn/cryptellation/internal/infra/indicators/db/sql/entities"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"gorm.io/gorm/clause"
)

func (d DB) GetSMA(ctx context.Context, payload db.ReadSMAPayload) (*timeserie.TimeSerie[float64], error) {
	tx := d.client.Where(`
		exchange_name = ? AND
		pair_symbol = ? AND
		period_symbol = ? AND
		period_number = ? AND
		price_type = ? AND
		time BETWEEN ? AND ?`,
		payload.ExchangeName,
		payload.PairSymbol,
		payload.Period.String(),
		payload.PeriodNumber,
		payload.PriceType.String(),
		payload.Start, payload.End)

	ent := []entities.SimpleMovingAverage{}
	if err := tx.WithContext(ctx).Find(&ent).Error; err != nil {
		return nil, err
	}

	return entities.FromEntityListToModelList(ent)
}

func (d DB) UpsertSMA(ctx context.Context, payload db.WriteSMAPayload) error {
	listCE := entities.FromModelListToEntityList(
		payload.ExchangeName,
		payload.PairSymbol,
		payload.Period.String(),
		int(payload.PeriodNumber),
		payload.PriceType,
		payload.TimeSerie,
	)
	return d.client.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&listCE).Error
}
