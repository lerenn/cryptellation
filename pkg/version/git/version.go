package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"golang.org/x/mod/semver"
)

var (
	// ErrNoVersion is the error encountered when no version is found in the tags.
	ErrNoVersion = fmt.Errorf("no version in the tags")
)

// ModuleVersions returns all the versions of a module in the git repository at path.
func ModuleVersions(path, appName string) ([]string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get repository versions
	tagsIter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	// Filter by service
	versions := make([]string, 0)
	if err := tagsIter.ForEach(func(r *plumbing.Reference) error {
		n, v := getNameVersionFromGitTagRef(r)
		if n == appName {
			versions = append(versions, v)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return versions, nil
}

// ModuleVersionsWithHash returns all the versions of a module in the git
// repository at path with their hash.
func ModuleVersionsWithHash(path, appName string) (map[string]string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get repository versions
	tagsIter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	// Filter by service
	versions := make(map[string]string, 0)
	if err := tagsIter.ForEach(func(r *plumbing.Reference) error {
		n, v := getNameVersionFromGitTagRef(r)
		if n == appName {
			versions[v] = r.Hash().String()
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return versions, nil
}

func getNameVersionFromGitTagRef(r *plumbing.Reference) (name string, version string) {
	parts := strings.Split(r.Name().Short(), "/")
	if len(parts) == 1 {
		version = parts[0]
	} else {
		name, version = strings.Join(parts[:len(parts)-1], "/"), parts[len(parts)-1]
	}

	return
}

// ModuleVersionsFromCurrentBranch returns all the version of a module in the
// git repository at path from the current branch.
func ModuleVersionsFromCurrentBranch(path, appName string) ([]string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get all versions from repository
	allVersions, err := ModuleVersionsWithHash(path, appName)
	if err != nil {
		return nil, err
	}

	// Get commit iterator
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	if err := commitIter.ForEach(func(r *object.Commit) error {
		for n, t := range allVersions {
			if t == r.Hash.String() {
				versions = append(versions, n)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return versions, nil
}

// LastModuleVersionFromCurrentBranch returns the last version of a module in
// the git repository at path from the current branch.
func LastModuleVersionFromCurrentBranch(path, appName string) (string, error) {
	// Get version from branch
	versions, err := ModuleVersionsFromCurrentBranch(path, appName)
	if err != nil {
		return "", err
	}

	// Check there is still versions
	if len(versions) == 0 {
		return "", ErrNoVersion
	}

	// Only get the list of versions
	semver.Sort(versions)
	return versions[len(versions)-1], nil
}
