package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

// Parameters is the entity for the parameters of a backtest.
type Parameters struct {
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time"`
	Mode        string    `bson:"mode"`
	PricePeriod string    `bson:"price_period"`
}

// Backtest is the entity for a backtest.
type Backtest struct {
	ID                string             `bson:"_id"`
	Parameters        Parameters         `bson:"parameters"`
	CurrentTime       time.Time          `bson:"current_time"`
	CurrentPriceType  string             `bson:"current_price_type"`
	Balances          []Balance          `bson:"balances"`
	Orders            []Order            `bson:"orders"`
	TickSubscriptions []TickSubscription `bson:"tick_subscriptions"`
}

// ToModel converts the entity to a model.
func (bt Backtest) ToModel() (backtest.Backtest, error) {
	priceType := candlestick.PriceType(bt.CurrentPriceType)
	if err := priceType.Validate(); err != nil {
		wrappedErr := fmt.Errorf("error when validating current price type, got %q: %w", bt.CurrentPriceType, err)
		return backtest.Backtest{}, wrappedErr
	}

	periodBetweenEvents := period.Symbol(bt.Parameters.PricePeriod)
	if err := periodBetweenEvents.Validate(); err != nil {
		return backtest.Backtest{}, err
	}

	mode := backtest.Mode(bt.Parameters.Mode)
	if err := mode.Validate(); err != nil {
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
		Parameters: backtest.Settings{
			StartTime:   bt.Parameters.StartTime,
			EndTime:     bt.Parameters.EndTime,
			Mode:        mode,
			PricePeriod: periodBetweenEvents,
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

// FromBacktestModel converts a model into an entity.
func FromBacktestModel(bt backtest.Backtest) Backtest {
	return Backtest{
		ID: bt.ID.String(),
		Parameters: Parameters{
			StartTime:   bt.Parameters.StartTime,
			EndTime:     bt.Parameters.EndTime,
			Mode:        bt.Parameters.Mode.String(),
			PricePeriod: bt.Parameters.PricePeriod.String(),
		},
		CurrentTime:       bt.CurrentCandlestick.Time,
		CurrentPriceType:  bt.CurrentCandlestick.Price.String(),
		Balances:          FromAccountModels(bt.Accounts),
		Orders:            FromOrderModels(bt.Orders),
		TickSubscriptions: FromTickSubscriptionModels(bt.PricesSubscriptions),
	}
}
