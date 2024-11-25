package utils

// ToReference return the value in argument as a reference.
func ToReference[T any](t T) *T {
	return &t
}

// FromReferenceOrDefault return the value from the reference or the default value
// from the type is the reference is nil.
func FromReferenceOrDefault[T any](t *T) T {
	if t == nil {
		return *new(T)
	}

	return *t
}
