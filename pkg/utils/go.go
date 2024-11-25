package utils

import "runtime"

// GoVersion returns the go version used to compile the library.
func GoVersion() string {
	return runtime.Version()[2:]
}
