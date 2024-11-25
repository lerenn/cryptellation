package utils

import "os"

// TemporaryEnvVar sets a temporary environment variable for the duration of the test.
func TemporaryEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
