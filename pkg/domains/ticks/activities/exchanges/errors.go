package exchanges

import "errors"

var (
	// ErrInexistantExchange is the error when the exchange does not exist.
	ErrInexistantExchange = errors.New("inexistant exchange")
)
