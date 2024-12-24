package order

import "errors"

var (
	// ErrInvalidType is the error for an invalid order type.
	ErrInvalidType = errors.New("invalid-order-type")
)

// Type is the order type.
type Type string

const (
	// TypeIsMarket is the market order type.
	TypeIsMarket Type = "market"
)

var (
	// Types is the list of order types.
	Types = []Type{
		TypeIsMarket,
	}
)

// Validate validates the order type.
func (t Type) Validate() error {
	for _, vt := range Types {
		if t == vt {
			return nil
		}
	}

	return ErrInvalidType
}

// String returns the string representation of the order type.
func (t Type) String() string {
	return string(t)
}
