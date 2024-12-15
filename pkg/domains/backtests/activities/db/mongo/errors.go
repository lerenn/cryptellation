package mongo

import "fmt"

var (
	// ErrNilID is returned when the ID is nil.
	ErrNilID = fmt.Errorf("ID is nil")
)
