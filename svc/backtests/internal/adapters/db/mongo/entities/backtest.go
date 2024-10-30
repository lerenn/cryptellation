package entities

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
)

type Parameters struct {
	StartTime time.Time `bson:"start_time"`
	EndTime   time.Time `bson:"end_time"`
	Period    string    `bson:"period"`
}

type Backtest struct {
	ID                string             `bson:"_id"`
	Parameters        Parameters         `bson:"parameters"`
	CurrentTime       time.Time          `bson:"current_time"`
	CurrentPriceType  string             `bson:"current_price_type"`
	Balances          []Balance          `bson:"balances"`
	Orders            []Order            `bson:"orders"`
	TickSubscriptions []TickSubscription `bson:"tick_subscriptions"`
}

func (bt Backtest) ToModel() (backtest.Backtest, error) {
	priceType := candlestick.Price(bt.CurrentPriceType)
	if err := priceType.Validate(); err != nil {
		wrappedErr := fmt.Errorf("error when validating current price type, got %q: %w", bt.CurrentPriceType, err)
		return backtest.Backtest{}, wrappedErr
	}

	periodBetweenEvents := period.Symbol(bt.Parameters.Period)
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
		ID: id,
		Parameters: backtest.Parameters{
			StartTime: bt.Parameters.StartTime,
			EndTime:   bt.Parameters.EndTime,
			Period:    periodBetweenEvents,
		},
		CurrentCandlestick: backtest.CurrentCandlestick{
			Time:  bt.CurrentTime,
			Price: priceType,
		},
		Accounts:            ToAccountModels(bt.Balances),
		Orders:              orders,
		PricesSubscriptions: ToTickSubscriptionModels(bt.TickSubscriptions),
	}, nil
}

func FromBacktestModel(bt backtest.Backtest) Backtest {
	return Backtest{
		ID: bt.ID.String(),
		Parameters: Parameters{
			StartTime: bt.Parameters.StartTime,
			EndTime:   bt.Parameters.EndTime,
			Period:    bt.Parameters.Period.String(),
		},
		CurrentTime:       bt.CurrentCandlestick.Time,
		CurrentPriceType:  bt.CurrentCandlestick.Price.String(),
		Balances:          FromAccountModels(bt.Accounts),
		Orders:            FromOrderModels(bt.Orders),
		TickSubscriptions: FromTickSubscriptionModels(bt.PricesSubscriptions),
	}
}
