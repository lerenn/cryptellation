package db

import "fmt"

var (
	// ErrNilID is returned when the ID is nil.
	ErrNilID = fmt.Errorf("ID is nil")
	// ErrNoDocument is returned when no document is found.
	ErrNoDocument = fmt.Errorf("no document found")
)
