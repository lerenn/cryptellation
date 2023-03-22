package entities

import (
	"time"

	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/pkg/types/period"
)

type Backtest struct {
	ID                  uint `gorm:"primaryKey"`
	StartTime           time.Time
	CurrentTime         time.Time
	CurrentPriceType    string
	EndTime             time.Time
	PeriodBetweenEvents string
	Balances            []Balance
	Orders              []Order
	TickSubscriptions   []TickSubscription
}

func (bt Backtest) ToModel() (domain.Backtest, error) {
	priceType := candlestick.PriceType(bt.CurrentPriceType)
	if err := priceType.Validate(); err != nil {
		return domain.Backtest{}, err
	}

	periodBetweenEvents := period.Symbol(bt.PeriodBetweenEvents)
	if err := periodBetweenEvents.Validate(); err != nil {
		return domain.Backtest{}, err
	}

	orders, err := ToOrderModels(bt.Orders)
	if err != nil {
		return domain.Backtest{}, err
	}

	return domain.Backtest{
		ID:        bt.ID,
		StartTime: bt.StartTime,
		CurrentCsTick: domain.CurrentCsTick{
			Time:      bt.CurrentTime,
			PriceType: priceType,
		},
		EndTime:             bt.EndTime,
		PeriodBetweenEvents: periodBetweenEvents,
		Accounts:            ToAccountModels(bt.Balances),
		Orders:              orders,
		TickSubscriptions:   ToTickSubscriptionModels(bt.TickSubscriptions),
	}, nil
}

func FromBacktestModel(bt domain.Backtest) Backtest {
	return Backtest{
		ID:                  bt.ID,
		StartTime:           bt.StartTime,
		CurrentTime:         bt.CurrentCsTick.Time,
		CurrentPriceType:    bt.CurrentCsTick.PriceType.String(),
		EndTime:             bt.EndTime,
		PeriodBetweenEvents: bt.PeriodBetweenEvents.String(),
		Balances:            FromAccountModels(bt.ID, bt.Accounts),
		Orders:              FromOrderModels(bt.ID, bt.Orders),
		TickSubscriptions:   FromTickSubscriptionModels(bt.ID, bt.TickSubscriptions),
	}
}
