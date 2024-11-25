package version

import (
	git "github.com/lerenn/cryptellation/v1/pkg/version/git"
)

// SetFromGit sets the commit hash and the version from the last git.
func SetFromGit(path, appName string) error {
	if err := SetVersionFromGit(path, appName); err != nil {
		return err
	}

	return SetCommitHashFromGit(path)
}

// SetCommitHashFromGit sets the commit hash from the last git.
func SetCommitHashFromGit(path string) error {
	commitHash, err := CommitHashFromGit(path)
	if err != nil {
		return err
	}

	globalCommitHash = commitHash
	return nil
}

// CommitHashFromGit returns the commit hash from the last git.
func CommitHashFromGit(path string) (string, error) {
	return git.LastCommitHash(path)
}

// SetVersionFromGit sets the version from the last git.
func SetVersionFromGit(path, appName string) error {
	version, err := FromGit(path, appName)
	if err != nil {
		return err
	}

	globalVersion = version
	return nil
}

// FromGit returns the version from the last git.
func FromGit(path, appName string) (string, error) {
	return git.LastModuleVersionFromCurrentBranch(path, appName)
}

// FullVersionFromGit returns the full version from the last git.
func FullVersionFromGit(path, appName string) (string, error) {
	version, err := FromGit(path, appName)
	if err != nil {
		return "", err
	}

	commitHash, err := CommitHashFromGit(path)
	if err != nil {
		return "", err
	}

	return fullVersion(version, commitHash), nil
}
