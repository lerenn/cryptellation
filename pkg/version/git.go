package version

import (
	"fmt"

	"golang.org/x/mod/semver"

	git "github.com/lerenn/cryptellation/pkg/vcs/git"
)

// SetFromGit sets the commit hash and the version from the last git
func SetFromGit(path, appName string) error {
	if err := SetVersionFromGit(path, appName); err != nil {
		return err
	}

	return SetCommitHashFromGit(path)
}

func SetCommitHashFromGit(path string) error {
	commitHash, err := CommitHashFromGit(path)
	if err != nil {
		return err
	}

	globalCommitHash = commitHash
	return nil
}

func CommitHashFromGit(path string) (string, error) {
	return git.LastCommitHash(path)
}

func SetVersionFromGit(path, appName string) error {
	version, err := VersionFromGit(path, appName)
	if err != nil {
		return err
	}

	globalVersion = version
	return nil
}

func VersionFromGit(path, appName string) (string, error) {
	// Get version from branch
	versions, err := git.GetServiceVersionsFromCurrentBranch(path, appName)
	if err != nil {
		return "", err
	}

	// Check there is still versions
	if len(versions) == 0 {
		return "", fmt.Errorf("no version in the tags")
	}

	// Only get the list of versions
	semver.Sort(versions)
	return versions[len(versions)-1], nil
}

func FullVersionFromGit(path, appName string) (string, error) {
	version, err := VersionFromGit(path, appName)
	if err != nil {
		return "", err
	}

	commitHash, err := CommitHashFromGit(path)
	if err != nil {
		return "", err
	}

	return fullVersion(version, commitHash), nil
}
