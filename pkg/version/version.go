package version

import (
	"fmt"
)

const (
	DefaultHash = "devel" // Default commit hash in case no value was provided to override.
)

type Version struct {
	// SemVer is the semantic version of the application
	SemVer string

	// Revision of the application
	// "dev" is the default hash is nothing is provided,
	// this ensure to show this was not built via a pipeline where the hash should be automatically passed
	CommitHash string
}

// returns a string representing the version and commit hash concatenated separated by a -
//
//	returns only the version if the commit hash is not defined
func (v Version) FullVersion() string {
	if v.CommitHash == "" {
		return v.SemVer
	}
	return fmt.Sprintf("%s-%s", v.SemVer, v.CommitHash)
}
