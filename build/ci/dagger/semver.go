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

func updateSemVer(major, minor, patch int, title string) (newMajor, newMinor, newPatch int) {
	switch {
	case strings.HasPrefix(title, "BREAKING CHANGE"):
		newMajor = major + 1
		newMinor = 0
		newPatch = 0
	case strings.HasPrefix(title, "feat"):
		newMajor = major
		newMinor = minor + 1
		newPatch = 0
	case strings.HasPrefix(title, "fix"):
		newMajor = major
		newMinor = minor
		newPatch = patch + 1
	default:
		newMajor = major
		newMinor = minor
		newPatch = patch
	}

	return newMajor, newMinor, newPatch
}

func parseAndUpdateSemVer(version, title string) (newVersion string, err error) {
	// Parse version
	major, minor, patch, err := parseSemVer(version)
	if err != nil {
		return "", err
	}

	// Update version
	newMajor, newMinor, newPatch := updateSemVer(major, minor, patch, title)

	// Return new version
	return fmt.Sprintf("%d.%d.%d", newMajor, newMinor, newPatch), nil
}
