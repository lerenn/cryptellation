package entities

import "strings"

func BinanceSymbol(pairSymbol string) string {
	return strings.ReplaceAll(pairSymbol, "-", "")
}
