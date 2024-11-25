package entities

import "strings"

// BinanceSymbol will convert a pair to a Binance symbol.
func BinanceSymbol(pair string) string {
	return strings.ReplaceAll(pair, "-", "")
}
