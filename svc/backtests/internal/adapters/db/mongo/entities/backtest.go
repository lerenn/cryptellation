package entities

import (
	"time"

	"cryptellation/svc/backtests/pkg/backtest"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
)

type Backtest struct {
	ID                  string             `bson:"_id"`
	StartTime           time.Time          `bson:"start_time"`
	CurrentTime         time.Time          `bson:"current_time"`
	CurrentPriceType    string             `bson:"current_price_type"`
	EndTime             time.Time          `bson:"end_time"`
	PeriodBetweenEvents string             `bson:"period_between_events"`
	Balances            []Balance          `bson:"balances"`
	Orders              []Order            `bson:"orders"`
	TickSubscriptions   []TickSubscription `bson:"tick_subscriptions"`
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

	id, err := uuid.Parse(bt.ID)
	if err != nil {
		return backtest.Backtest{}, err
	}

	return backtest.Backtest{
		ID:        id,
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
		ID:                  bt.ID.String(),
		StartTime:           bt.StartTime,
		CurrentTime:         bt.CurrentCsTick.Time,
		CurrentPriceType:    bt.CurrentCsTick.PriceType.String(),
		EndTime:             bt.EndTime,
		PeriodBetweenEvents: bt.PeriodBetweenEvents.String(),
		Balances:            FromAccountModels(bt.Accounts),
		Orders:              FromOrderModels(bt.Orders),
		TickSubscriptions:   FromTickSubscriptionModels(bt.TickSubscriptions),
	}
}
