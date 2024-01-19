package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const defaultTestVersion = "1.2.3"

func TestGetVersion(t *testing.T) {
	testCases := []struct {
		name          string
		version       string
		expectedValue string
	}{
		{"2dotted", "1.2.3", "1.2.3"},
		{"1dotted", "1.23", "1.23"},
		{"empty", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			globalVersion = tc.version
			assert.Equal(t, tc.expectedValue, Version())
		})
	}
}

func TestSetVersion(t *testing.T) {
	testCases := []struct {
		name          string
		version       string
		expectedValue string
	}{
		{"validVersion", "1.2.3", "1.2.3"},
		{"empty", "", defaultTestVersion},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			globalVersion = defaultTestVersion
			SetVersion(tc.version)
			assert.Equal(t, tc.expectedValue, globalVersion)
		})
	}
}

func TestSetCommitHash(t *testing.T) {
	testCases := []struct {
		name          string
		commitHash    string
		expectedValue string
	}{
		{"validVersion", "abcdefgh12345", "abcdefgh12345"},
		{"empty", "", ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			globalCommitHash = DefaultHash
			SetCommitHash(tc.commitHash)
			assert.Equal(t, tc.expectedValue, globalCommitHash)
		})
	}
}

func TestGetCommitHash(t *testing.T) {
	testCases := []struct {
		name          string
		commitHash    string
		expectedValue string
	}{
		{"validVersion", "abcdefgh12345", "abcdefgh12345"},
		{"DefaultHash", DefaultHash, DefaultHash},
		{"empty", "", ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			globalCommitHash = DefaultHash
			SetCommitHash(tc.commitHash)
			assert.Equal(t, tc.expectedValue, CommitHash())
		})
	}
}

func TestGetFullVersion(t *testing.T) {
	// Test case setup
	testCases := []struct {
		name                string
		version, commitHash string
		expectedValue       string
	}{
		{"noUpdate", "", "", "1.2.3"},
		{"versionOnly", "1.23", "", "1.23"},
		{"commitHashOnly", "", "abcef835", "1.2.3-abcef835"},
		{"bothUpdated", "0.0.3", "abcef835", "0.0.3-abcef835"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			globalVersion, globalCommitHash = defaultTestVersion, "dev" // Setup default values before each run
			SetVersion(tc.version)
			SetCommitHash(tc.commitHash)
			assert.Equal(t, tc.expectedValue, FullVersion())
		})
	}
}
