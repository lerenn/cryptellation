package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseSemVer(version string) (major, minor, patch int, err error) {
	// Remove wrong characters
	version = strings.TrimPrefix(version, "v")
	version = strings.Trim(version, "\"")

	// Split version into parts
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, errors.New("invalid version format:" + version)
	}

	// Convert parts to integers
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

type SemVerChange string

const (
	SemVerChangeMajor   SemVerChange = "major"
	SemVerChangeMinor   SemVerChange = "minor"
	SemVerChangePatch   SemVerChange = "patch"
	SemVerChangeNone    SemVerChange = "none"
	SemVerChangeUnknown SemVerChange = ""
)

func processSemVerChange(version, title string) (change SemVerChange, newVersion string, err error) {
	// Parse version
	major, minor, patch, err := parseSemVer(version)
	if err != nil {
		return SemVerChangeUnknown, "", err
	}

	// Get semver change
	var semVerChange SemVerChange
	switch {
	case strings.HasPrefix(title, "BREAKING CHANGE"):
		semVerChange = SemVerChangeMajor
	case strings.HasPrefix(title, "feat"):
		semVerChange = SemVerChangeMinor
	case strings.HasPrefix(title, "fix"):
		semVerChange = SemVerChangePatch
	default:
		semVerChange = SemVerChangeNone
	}

	// Update version
	switch semVerChange {
	case SemVerChangeMajor:
		major = major + 1
		minor = 0
		patch = 0
	case SemVerChangeMinor:
		minor = minor + 1
		patch = 0
	case SemVerChangePatch:
		patch = patch + 1
	default:
		// No change
	}

	// Return new version
	return semVerChange, fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}
