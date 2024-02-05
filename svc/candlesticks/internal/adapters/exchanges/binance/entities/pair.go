package entities

import "strings"

func BinanceSymbol(pair string) string {
	return strings.ReplaceAll(pair, "-", "")
}
