package main

func toReference[T any](t T) *T {
	return &t
}
