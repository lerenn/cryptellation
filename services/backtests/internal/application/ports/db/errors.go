package db

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNotImplemented = errors.New("not implemented")
)
