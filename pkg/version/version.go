package version

import (
	"fmt"
)

const (
	DefaultHash = "dev" // Default commit hash in case no value was provided to override.
)

var (
	// Version of the application
	version = "1.0.0"

	// Revision of the application
	// "dev" is the default hash is nothing is provided,
	// this ensure to show this was not built via a pipeline where the hash should be automatically passed
	commitHash = DefaultHash
)

// returns a string representing the current version
func GetVersion() string {
	return version
}

// sets the version to tha value provided as ver unless ver is empty
func SetVersion(ver string) {
	if ver != "" {
		version = ver
	}
}

// sets the commit hash of the application to the value provided as hash
// empty values are accepted
func SetCommitHash(hash string) {
	commitHash = hash
}

// returns a string representing the current commitHash
func GetCommitHash() string {
	return commitHash
}

// returns a string representing the version and commit hash concatenated separated by a -
//
//	returns only the version if the commit hash is not defined
func GetFullVersion() string {
	if commitHash == "" {
		return GetVersion()
	}
	return fmt.Sprintf("%s-%s", version, commitHash)
}
