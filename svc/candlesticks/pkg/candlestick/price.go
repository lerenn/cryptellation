package candlestick

import "errors"

var (
	ErrInvalidPriceType = errors.New("invalid-price-type")
)

type Price string

const (
	PriceIsOpen  Price = "open"
	PriceIsHigh  Price = "high"
	PriceIsLow   Price = "low"
	PriceIsClose Price = "close"
)

var Prices = []Price{
	PriceIsOpen,
	PriceIsHigh,
	PriceIsLow,
	PriceIsClose,
}

func (pt Price) String() string {
	return string(pt)
}

func (pt Price) Validate() error {
	for _, vpt := range Prices {
		if vpt.String() == pt.String() {
			return nil
		}
	}

	return ErrInvalidPriceType
}
