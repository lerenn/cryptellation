package forwardtest

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

var (
	// ErrEmptyAccounts is returned when the accounts are empty.
	ErrEmptyAccounts = errors.New("empty accounts")
	// ErrInvalidExchange is returned when the exchange is invalid.
	ErrInvalidExchange = errors.New("invalid exchange")
)

// Forwardtest is a forwardtest.
type Forwardtest struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Accounts  map[string]account.Account
	Orders    []order.Order
}

// NewForwardtestParams is the params for the New function.
type NewForwardtestParams struct {
	Accounts map[string]account.Account
}

// Validate validates the NewParams.
func (np NewForwardtestParams) Validate() error {
	if len(np.Accounts) == 0 {
		return ErrEmptyAccounts
	}

	return nil
}

// New creates a new forwardtest.
func New(params NewForwardtestParams) Forwardtest {
	return Forwardtest{
		ID:       uuid.New(),
		Accounts: params.Accounts,
	}
}

// AddOrder adds an order to the forwardtest.
func (ft *Forwardtest) AddOrder(o order.Order, cs candlestick.Candlestick) error {
	// Get exchange account
	exchangeAccount, ok := ft.Accounts[o.Exchange]
	if !ok {
		return fmt.Errorf("error with orders exchange %q: %w", o.Exchange, ErrInvalidExchange)
	}

	// Get price
	price := cs.Close
	if price == 0 {
		return errors.New("price is 0, that should not happen")
	}

	// Apply order
	if err := exchangeAccount.ApplyOrder(price, o); err != nil {
		return err
	}
	ft.Accounts[o.Exchange] = exchangeAccount

	// Update and save the order
	o.ExecutionTime = utils.ToReference(time.Now())
	o.Price = price
	ft.Orders = append(ft.Orders, o)

	return nil
}

// GetAccountsSymbols returns the list of symbols used in the accounts.
func (ft Forwardtest) GetAccountsSymbols() []string {
	symbols := make(map[string]string, 0)

	for _, account := range ft.Accounts {
		for symbol := range account.Balances {
			if _, ok := symbols[symbol]; !ok {
				symbols[symbol] = symbol
			}
		}
	}

	return utils.MapToList(symbols)
}
