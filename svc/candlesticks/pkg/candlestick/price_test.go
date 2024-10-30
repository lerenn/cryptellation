package candlestick

import "testing"

func TestCandlestickPrices(t *testing.T) {
	validPrices := []string{
		"open",
		"high",
		"low",
		"close",
	}

	for _, vpt := range validPrices {
		if err := Price(vpt).Validate(); err != nil {
			t.Error("Price should be valid: ", vpt)
		}
	}

	if err := Price("unknown").Validate(); err == nil {
		t.Error("Price should not be valid: unknown")
	}
}
