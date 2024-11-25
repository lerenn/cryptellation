package candlestick

import "errors"

var (
	// ErrInvalidPriceType is returned when the price type is invalid.
	ErrInvalidPriceType = errors.New("invalid-price-type")
)

// Price is the type of price to use.
type Price string

const (
	// PriceIsOpen is the open price.
	PriceIsOpen Price = "open"
	// PriceIsHigh is the high price.
	PriceIsHigh Price = "high"
	// PriceIsLow is the low price.
	PriceIsLow Price = "low"
	// PriceIsClose is the close price.
	PriceIsClose Price = "close"
)

// Prices is the list of all available prices.
var Prices = []Price{
	PriceIsOpen,
	PriceIsHigh,
	PriceIsLow,
	PriceIsClose,
}

// String returns the string representation of the price type.
func (pt Price) String() string {
	return string(pt)
}

// Validate checks if the price type is valid.
func (pt Price) Validate() error {
	for _, vpt := range Prices {
		if vpt.String() == pt.String() {
			return nil
		}
	}

	return ErrInvalidPriceType
}
