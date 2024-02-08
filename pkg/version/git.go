package version

import (
	git "github.com/lerenn/cryptellation/pkg/version/git"
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
	return git.LastModuleVersionFromCurrentBranch(path, appName)
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
