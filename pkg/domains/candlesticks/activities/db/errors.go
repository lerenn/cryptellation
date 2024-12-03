package db

import "errors"

var (
	// ErrNotFound is returned when the requested data is not found.
	ErrNotFound = errors.New("not-found")
)
