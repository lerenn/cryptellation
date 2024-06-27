package forwardtest

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

var (
	ErrEmptyAccounts   = errors.New("empty accounts")
	ErrInvalidExchange = errors.New("invalid exchange")
)

type ForwardTest struct {
	ID       uuid.UUID
	Accounts map[string]account.Account
	Orders   []order.Order
}

type NewPayload struct {
	Accounts map[string]account.Account
}

func (np NewPayload) Validate() error {
	if len(np.Accounts) == 0 {
		return ErrEmptyAccounts
	}

	return nil
}

func New(payload NewPayload) ForwardTest {
	return ForwardTest{
		ID:       uuid.New(),
		Accounts: payload.Accounts,
	}
}

func (ft *ForwardTest) AddOrder(o order.Order, cs candlestick.Candlestick) error {
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

func (ft ForwardTest) GetAccountsSymbols() []string {
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
