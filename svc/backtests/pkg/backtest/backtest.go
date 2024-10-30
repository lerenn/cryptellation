package backtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/pkg/utils"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/google/uuid"
)

var (
	ErrTickSubscriptionAlreadyExists = errors.New("tick subscription already exists")
	ErrInvalidExchange               = errors.New("invalid exchange")
	ErrNoDataForOrderValidation      = errors.New("no data for order validation")
	ErrStartAfterEnd                 = errors.New("start after end")
	ErrInvalidPricePeriod            = errors.New("invalid price period")
)

// Current candlestick based on candlestick step
type CurrentCandlestick struct {
	Time  time.Time
	Price candlestick.Price
}

type Parameters struct {
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Mode        Mode          `json:"mode"`
	PricePeriod period.Symbol `json:"price_period"`
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
	Accounts  map[string]account.Account
	StartTime time.Time
	EndTime   *time.Time

	Mode        *Mode
	PricePeriod *period.Symbol
}

func (payload *NewPayload) EmptyFieldsToDefault() *NewPayload {
	if payload.EndTime == nil {
		payload.EndTime = defaultEndTime()
	}

	if payload.PricePeriod == nil {
		payload.PricePeriod = period.M1.Opt()
	}

	if payload.Mode == nil {
		payload.Mode = utils.ToReference(ModeIsCloseOHLC)
	}

	return payload
}

func (payload NewPayload) Validate() error {
	if !payload.StartTime.Before(*payload.EndTime) {
		return ErrStartAfterEnd
	}

	if payload.PricePeriod == nil {
		return fmt.Errorf("%w: nil", ErrInvalidPricePeriod)
	}

	if payload.Mode == nil {
		return ErrInvalidMode
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

func New(payload NewPayload) (Backtest, error) {
	// Set default fields payload and validate it
	if err := payload.EmptyFieldsToDefault().Validate(); err != nil {
		return Backtest{}, err
	}

	// Set current candlestick based on mode
	cc := CurrentCandlestick{
		Time: payload.StartTime,
	}
	switch *payload.Mode {
	case ModeIsCloseOHLC:
		cc.Price = candlestick.PriceIsClose
	case ModeIsFullOHLC:
		cc.Price = candlestick.PriceIsOpen
	}

	return Backtest{
		ID: uuid.New(),
		Parameters: Parameters{
			StartTime:   payload.StartTime,
			EndTime:     *payload.EndTime,
			Mode:        *payload.Mode,
			PricePeriod: *payload.PricePeriod,
		},
		CurrentCandlestick:  cc,
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

func (bt *Backtest) Advance() (done bool, err error) {
	switch bt.Parameters.Mode {
	case ModeIsCloseOHLC:
		bt.advanceWithModeIsCloseOHLC()
	case ModeIsFullOHLC:
		bt.advanceWithModeIsFullOHLC()
	default:
		return false, fmt.Errorf("error with backtest mode %q: %w", bt.Parameters.Mode, ErrInvalidMode)
	}

	return bt.Done(), nil
}

func (bt *Backtest) advanceWithModeIsCloseOHLC() {
	bt.SetCurrentTime(bt.CurrentCandlestick.Time.Add(bt.Parameters.PricePeriod.Duration()))
}

func (bt *Backtest) advanceWithModeIsFullOHLC() {
	switch bt.CurrentCandlestick.Price {
	case candlestick.PriceIsOpen:
		bt.CurrentCandlestick.Price = candlestick.PriceIsHigh
	case candlestick.PriceIsHigh:
		bt.CurrentCandlestick.Price = candlestick.PriceIsLow
	case candlestick.PriceIsLow:
		bt.CurrentCandlestick.Price = candlestick.PriceIsClose
	case candlestick.PriceIsClose:
		bt.SetCurrentTime(bt.CurrentCandlestick.Time.Add(bt.Parameters.PricePeriod.Duration()))
	default:
		bt.CurrentCandlestick.Price = candlestick.PriceIsOpen
	}
}

func (bt Backtest) Done() bool {
	return !bt.CurrentCandlestick.Time.Before(bt.Parameters.EndTime)
}

func (bt *Backtest) SetCurrentTime(ts time.Time) {
	// Set new time
	bt.CurrentCandlestick.Time = ts

	// Starting the time on open if mode is full OHLC
	if bt.Parameters.Mode == ModeIsFullOHLC {
		bt.CurrentCandlestick.Price = candlestick.PriceIsOpen
	}
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
