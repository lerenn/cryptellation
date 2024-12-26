package candlestick

import "errors"

var (
	// ErrInvalidPriceType is returned when the price type is invalid.
	ErrInvalidPriceType = errors.New("invalid-price-type")
)

// PriceType is the type of price to use.
type PriceType string

const (
	// PriceTypeIsOpen is the open price.
	PriceTypeIsOpen PriceType = "open"
	// PriceTypeIsHigh is the high price.
	PriceTypeIsHigh PriceType = "high"
	// PriceTypeIsLow is the low price.
	PriceTypeIsLow PriceType = "low"
	// PriceTypeIsClose is the close price.
	PriceTypeIsClose PriceType = "close"
)

// PriceTypes is the list of all available price types.
var PriceTypes = []PriceType{
	PriceTypeIsOpen,
	PriceTypeIsHigh,
	PriceTypeIsLow,
	PriceTypeIsClose,
}

// String returns the string representation of the price type.
func (pt PriceType) String() string {
	return string(pt)
}

// Validate checks if the price type is valid.
func (pt PriceType) Validate() error {
	for _, vpt := range PriceTypes {
		if vpt.String() == pt.String() {
			return nil
		}
	}

	return ErrInvalidPriceType
}
