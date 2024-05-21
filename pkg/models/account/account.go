package account

import (
	"errors"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/pair"
)

var (
	ErrInvalidBalanceAmount = errors.New("invalid balance amount")
	ErrInvalidBalanceAsset  = errors.New("invalid balance asset")
	ErrNotEnoughAsset       = errors.New("not enough asset")
)

type Account struct {
	Balances map[string]float64
}

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

func (a *Account) ApplyOrder(price float64, o order.Order) error {
	// Get base and quote based on symbol
	baseSymbol, quoteSymbol, err := pair.ParsePair(o.Pair)
	if err != nil {
		return fmt.Errorf("error when parsing order pair symbol: %w", err)
	}

	// Apply order
	quoteEquivalentQty := price * o.Quantity
	if o.Side == order.SideIsBuy {
		available, ok := a.Balances[quoteSymbol]
		if !ok {
			return fmt.Errorf("%w: no %s on %s", ErrNotEnoughAsset, quoteSymbol, o.Pair)
		} else if quoteEquivalentQty > available {
			return fmt.Errorf(
				"%w: not enough %s on %s (min=%f, got=%f)",
				ErrNotEnoughAsset, quoteSymbol, o.Pair,
				quoteEquivalentQty, available)
		}

		a.Balances[quoteSymbol] -= quoteEquivalentQty
		a.Balances[baseSymbol] += o.Quantity
	} else {
		available, ok := a.Balances[baseSymbol]
		if !ok {
			return fmt.Errorf("%w: no %s on %s", ErrNotEnoughAsset, baseSymbol, o.Pair)
		} else if o.Quantity > available {
			return fmt.Errorf(
				"%w: not enough %s on %s (min=%f, got=%f)",
				ErrNotEnoughAsset, baseSymbol, o.Pair,
				o.Quantity, available)
		}

		a.Balances[quoteSymbol] += quoteEquivalentQty
		a.Balances[baseSymbol] -= o.Quantity
	}

	return nil
}
