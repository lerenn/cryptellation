package order

import (
	"errors"
)

var (
	// ErrInvalidSide is the error for an invalid order side.
	ErrInvalidSide = errors.New("invalid-order-side")
)

// Side is the side of an order.
type Side string

const (
	// SideIsBuy is the buy side of an order.
	SideIsBuy Side = "buy"
	// SideIsSell is the sell side of an order.
	SideIsSell Side = "sell"
)

var (
	// Sides is the list of all the possible sides for an order.
	Sides = []Side{
		SideIsBuy,
		SideIsSell,
	}
)

// Validate validates the order side.
func (s Side) Validate() error {
	for _, vs := range Sides {
		if s == vs {
			return nil
		}
	}

	return ErrInvalidSide
}

// String returns the string representation of the order side.
func (s Side) String() string {
	return string(s)
}
