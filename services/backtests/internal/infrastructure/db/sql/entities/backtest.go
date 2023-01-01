package entities

import (
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
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

func (bt Backtest) ToModel() (backtest.Backtest, error) {
	priceType := candlestick.PriceType(bt.CurrentPriceType)
	if err := priceType.Validate(); err != nil {
		return backtest.Backtest{}, err
	}

	periodBetweenEvents := period.Symbol(bt.PeriodBetweenEvents)
	if err := periodBetweenEvents.Validate(); err != nil {
		return backtest.Backtest{}, err
	}

	orders, err := ToOrderModels(bt.Orders)
	if err != nil {
		return backtest.Backtest{}, err
	}

	return backtest.Backtest{
		ID:        bt.ID,
		StartTime: bt.StartTime,
		CurrentCsTick: backtest.CurrentCsTick{
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

func FromBacktestModel(bt backtest.Backtest) Backtest {
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
