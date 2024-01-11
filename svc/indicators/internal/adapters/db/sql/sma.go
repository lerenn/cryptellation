package sql

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/timeserie"
	"github.com/lerenn/cryptellation/svc/indicators/internal/adapters/db/sql/entities"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
	"gorm.io/gorm/clause"
)

func (a Adapter) GetSMA(ctx context.Context, payload db.ReadSMAPayload) (*timeserie.TimeSerie[float64], error) {
	tx := a.db.Client.Where(`
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

func (a Adapter) UpsertSMA(ctx context.Context, payload db.WriteSMAPayload) error {
	listCE := entities.FromModelListToEntityList(
		payload.ExchangeName,
		payload.PairSymbol,
		payload.Period.String(),
		int(payload.PeriodNumber),
		payload.PriceType,
		payload.TimeSerie,
	)
	return a.db.Client.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&listCE).Error
}
