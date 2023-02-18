package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullVersion(t *testing.T) {
	// Test case setup
	testCases := []struct {
		name               string
		semVer, commitHash string
		expectedValue      string
	}{
		{"semVerOnly", "1.23", "", "1.23"},
		{"everything", "0.0.3", "abcef835", "0.0.3-abcef835"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := Version{SemVer: tc.semVer, CommitHash: tc.commitHash}
			assert.Equal(t, tc.expectedValue, v.FullVersion(), tc.name)
		})
	}
}
