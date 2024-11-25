package pair

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidPair is returned when the pair symbol is invalid.
	ErrInvalidPair = errors.New("invalid pair symbol")
)

// FormatPair formats a pair symbol.
func FormatPair(baseSymbol, quoteSymbol string) string {
	return fmt.Sprintf("%s-%s", baseSymbol, quoteSymbol)
}

// ParsePair parses a pair symbol.
func ParsePair(symbol string) (baseSymbol, quoteSymbol string, err error) {
	split := strings.Split(symbol, "-")
	if len(split) != 2 {
		return "", "", fmt.Errorf("error parsing pair symbol %q: %w", symbol, ErrInvalidPair)
	}

	return split[0], split[1], nil
}
