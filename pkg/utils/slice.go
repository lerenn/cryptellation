package utils

// MergeSliceIntoUnique merges two slices into one, keeping only unique values.
func MergeSliceIntoUnique[T comparable](slice1, slice2 []T) []T {
	// Create a map of the first array
	m1 := make(map[T]bool, len(slice1))
	for _, s1 := range slice1 {
		m1[s1] = true
	}

	// Create a map of the second array
	m2 := make(map[T]bool, len(slice2))
	for _, s2 := range slice2 {
		m2[s2] = true
	}

	// Merge the two maps
	for k := range m2 {
		m1[k] = true
	}

	// Create the resulting array
	res := make([]T, 0, len(m1))
	for k := range m1 {
		res = append(res, k)
	}

	return res
}
