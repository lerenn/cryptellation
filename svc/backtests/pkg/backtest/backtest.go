package backtest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
)

var (
	ErrTickSubscriptionAlreadyExists = errors.New("tick subscription already exists")
	ErrInvalidExchange               = errors.New("invalid exchange")
	ErrNoDataForOrderValidation      = errors.New("no data for order validation")
	ErrStartAfterEnd                 = errors.New("start after end")
)

// Current candlestick based on candlestick step
type CurrentCandlestick struct {
	Time  time.Time
	Price candlestick.Price
}

type Parameters struct {
	StartTime time.Time
	EndTime   time.Time
	Mode      Mode
	// Period between events (only in OHLC modes)
	Period period.Symbol `json:"period"`
}

type Backtest struct {
	ID                  uuid.UUID                  `json:"id"`
	Parameters          Parameters                 `json:"parameters"`
	CurrentCandlestick  CurrentCandlestick         `json:"current_candlestick"`
	Accounts            map[string]account.Account `json:"accounts"`
	PricesSubscriptions []event.PricesSubscription `json:"prices_subscriptions"`
	Orders              []order.Order              `json:"orders"`
}

type NewPayload struct {
	Accounts              map[string]account.Account
	StartTime             time.Time
	EndTime               *time.Time
	DurationBetweenEvents *time.Duration
}

func (payload *NewPayload) EmptyFieldsToDefault() *NewPayload {
	if payload.EndTime == nil {
		payload.EndTime = defaultEndTime()
	}

	if payload.DurationBetweenEvents == nil {
		d := time.Minute
		payload.DurationBetweenEvents = &d
	}

	return payload
}

func (payload NewPayload) Validate() error {
	if !payload.StartTime.Before(*payload.EndTime) {
		return ErrStartAfterEnd
	}

	for exchange, a := range payload.Accounts {
		if exchange == "" {
			return fmt.Errorf("error with exchange %q in new backtest payload: %w", exchange, ErrInvalidExchange)
		}

		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func defaultEndTime() *time.Time {
	t := time.Now()
	return &t
}

func New(ctx context.Context, payload NewPayload) (Backtest, error) {
	if err := payload.EmptyFieldsToDefault().Validate(); err != nil {
		return Backtest{}, err
	}

	per, err := period.FromDuration(*payload.DurationBetweenEvents)
	if err != nil {
		return Backtest{}, fmt.Errorf("invalid duration between candlesticks: %w", err)
	}

	return Backtest{
		ID: uuid.New(),
		Parameters: Parameters{
			StartTime: payload.StartTime,
			EndTime:   *payload.EndTime,
			Period:    per,
		},
		CurrentCandlestick: CurrentCandlestick{
			Time:  payload.StartTime,
			Price: candlestick.PriceIsOpen,
		},
		Accounts:            payload.Accounts,
		PricesSubscriptions: make([]event.PricesSubscription, 0),
		Orders:              make([]order.Order, 0),
	}, nil
}

func (bt Backtest) CurrentTime() string {
	return fmt.Sprintf("%s [%s]", bt.CurrentCandlestick.Time, bt.CurrentCandlestick.Price)
}

func (bt Backtest) MarshalBinary() ([]byte, error) {
	return json.Marshal(bt)
}

func (bt *Backtest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, bt)
}

func (bt *Backtest) Advance() (done bool) {
	return bt.advanceThroughTicks()
}

func (bt *Backtest) advanceThroughTicks() (done bool) {
	switch bt.CurrentCandlestick.Price {
	case candlestick.PriceIsOpen:
		bt.CurrentCandlestick.Price = candlestick.PriceIsHigh
	case candlestick.PriceIsHigh:
		bt.CurrentCandlestick.Price = candlestick.PriceIsLow
	case candlestick.PriceIsLow:
		bt.CurrentCandlestick.Price = candlestick.PriceIsClose
	case candlestick.PriceIsClose:
		bt.SetCurrentTime(bt.CurrentCandlestick.Time.Add(bt.Parameters.Period.Duration()))
	default:
		bt.CurrentCandlestick.Price = candlestick.PriceIsOpen
	}

	return bt.Done()
}

func (bt Backtest) Done() bool {
	return !bt.CurrentCandlestick.Time.Before(bt.Parameters.EndTime)
}

func (bt *Backtest) SetCurrentTime(ts time.Time) {
	bt.CurrentCandlestick.Time = ts
	bt.CurrentCandlestick.Price = candlestick.PriceIsOpen
}

func (bt *Backtest) CreateTickSubscription(exchange string, pair string) (event.PricesSubscription, error) {
	for _, ts := range bt.PricesSubscriptions {
		if ts.Exchange == exchange && ts.Pair == pair {
			return event.PricesSubscription{}, ErrTickSubscriptionAlreadyExists
		}
	}

	s := event.PricesSubscription{
		Exchange: exchange,
		Pair:     pair,
	}
	bt.PricesSubscriptions = append(bt.PricesSubscriptions, s)

	return s, nil
}

func (bt *Backtest) AddOrder(ord order.Order, cs candlestick.Candlestick) error {
	// Get exchange account
	exchangeAccount, ok := bt.Accounts[ord.Exchange]
	if !ok {
		return fmt.Errorf("error with orders exchange %q: %w", ord.Exchange, ErrInvalidExchange)
	}

	// Execute the order
	price := cs.Price(bt.CurrentCandlestick.Price)
	if err := exchangeAccount.ApplyOrder(price, ord); err != nil {
		return err
	}
	bt.Accounts[ord.Exchange] = exchangeAccount

	// Update and save the order
	ord.ExecutionTime = &bt.CurrentCandlestick.Time
	ord.Price = price
	bt.Orders = append(bt.Orders, ord)

	return nil
}
