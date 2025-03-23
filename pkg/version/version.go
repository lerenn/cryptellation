package version

import (
	"fmt"
)

const (
	// DefaultHash is the default hash if there is no hash provided.
	DefaultHash = "dev"
	// DefaultVersion is the default version if there is no hash provided.
	DefaultVersion = "1.12.0"
)

var (
	// Version of the application.
	globalVersion = DefaultVersion

	// Revision of the application.
	globalCommitHash = DefaultHash
)

// Version returns a string representing the current version.
func Version() string {
	return globalVersion
}

// SetVersion sets the version to tha value provided as ver unless ver is empty.
func SetVersion(ver string) {
	if ver != "" {
		globalVersion = ver
	}
}

// SetCommitHash sets the commit hash of the application to the value provided as hash.
// Empty values are accepted.
func SetCommitHash(hash string) {
	globalCommitHash = hash
}

// CommitHash returns a string representing the current commitHash.
func CommitHash() string {
	return globalCommitHash
}

// FullVersion returns a string representing the version and commit hash concatenated separated by a '-'.
//
// Returns only the version if the commit hash is not defined.
func FullVersion() string {
	if globalCommitHash == "" {
		return Version()
	}
	return fullVersion(globalVersion, globalCommitHash)
}

func fullVersion(version, commitHash string) string {
	return fmt.Sprintf("%s-%s", version, commitHash)
}
