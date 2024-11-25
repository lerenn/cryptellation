package db

import "errors"

var (
	// ErrNotFound is returned when the document is not found.
	ErrNotFound = errors.New("not-found")
)
