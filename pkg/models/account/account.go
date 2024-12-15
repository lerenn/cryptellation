package account

import (
	"errors"
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/models/order"
	"github.com/lerenn/cryptellation/v1/pkg/models/pair"
)

var (
	// ErrInvalidBalanceAmount is the error for an invalid balance amount.
	ErrInvalidBalanceAmount = errors.New("invalid balance amount")
	// ErrInvalidBalanceAsset is the error for an invalid balance asset.
	ErrInvalidBalanceAsset = errors.New("invalid balance asset")
	// ErrNotEnoughAsset is the error for not enough asset.
	ErrNotEnoughAsset = errors.New("not enough asset")
)

// Account is the struct for an account.
type Account struct {
	Balances map[string]float64 `json:"balances"`
}

// Validate will validate the account.
func (a Account) Validate() error {
	for asset, balance := range a.Balances {
		if asset == "" {
			return ErrInvalidBalanceAsset
		}

		if balance < 0 {
			return ErrInvalidBalanceAmount
		}
	}

	return nil
}

// ApplyOrder will apply an order to the account.
func (a *Account) ApplyOrder(price float64, o order.Order) error {
	switch o.Side {
	case order.SideIsBuy:
		return a.applyBuyOrder(price, o)
	case order.SideIsSell:
		return a.applySellOrder(price, o)
	default:
		return fmt.Errorf("unknown order side: %s", o.Side)
	}
}

func (a *Account) applyBuyOrder(price float64, o order.Order) error {
	// Get base and quote based on symbol
	baseSymbol, quoteSymbol, err := pair.ParsePair(o.Pair)
	if err != nil {
		return fmt.Errorf("error when parsing order pair symbol: %w", err)
	}

	// Check order
	quoteEquivalentQty := price * o.Quantity
	available, ok := a.Balances[quoteSymbol]
	if !ok {
		return fmt.Errorf("%w: no %s on %s", ErrNotEnoughAsset, quoteSymbol, o.Pair)
	} else if quoteEquivalentQty > available {
		return fmt.Errorf(
			"%w: not enough %s on %s (min=%f, got=%f)",
			ErrNotEnoughAsset, quoteSymbol, o.Pair,
			quoteEquivalentQty, available)
	}

	// Apply order on balances
	a.Balances[quoteSymbol] -= quoteEquivalentQty
	a.Balances[baseSymbol] += o.Quantity

	return nil
}

func (a *Account) applySellOrder(price float64, o order.Order) error {
	// Get base and quote based on symbol
	baseSymbol, quoteSymbol, err := pair.ParsePair(o.Pair)
	if err != nil {
		return fmt.Errorf("error when parsing order pair symbol: %w", err)
	}

	// Check order
	quoteEquivalentQty := price * o.Quantity
	available, ok := a.Balances[baseSymbol]
	if !ok {
		return fmt.Errorf("%w: no %s on %s", ErrNotEnoughAsset, baseSymbol, o.Pair)
	} else if o.Quantity > available {
		return fmt.Errorf(
			"%w: not enough %s on %s (min=%f, got=%f)",
			ErrNotEnoughAsset, baseSymbol, o.Pair,
			o.Quantity, available)
	}

	// Apply order on balances
	a.Balances[quoteSymbol] += quoteEquivalentQty
	a.Balances[baseSymbol] -= o.Quantity

	return nil
}
