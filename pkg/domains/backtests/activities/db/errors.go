package db

import (
	"errors"
)

var (
	// ErrNilID is returned when the ID is nil.
	ErrNilID = errors.New("ID is nil")
	// ErrNotFound is returned when the record is not found.
	ErrNotFound = errors.New("not found")
	// ErrNotImplemented is returned when the method is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)
