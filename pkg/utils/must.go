package utils

// Must panics if the error is not nil.
func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}

	return a
}
