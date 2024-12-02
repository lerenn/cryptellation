package exchanges

import "errors"

var (
	// ErrInexistantExchange is returned when the exchange does not exist.
	ErrInexistantExchange = errors.New("inexistant exchange")
)
