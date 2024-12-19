package db

import (
	"errors"
)

var (
	// ErrNilID is returned when the ID is nil.
	ErrNilID = errors.New("ID is nil")
	// ErrRecordNotFound is returned when the record is not found.
	ErrRecordNotFound = errors.New("record not found")
	// ErrNotImplemented is returned when the method is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)
