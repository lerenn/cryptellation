package git

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetServiceVersions(path, appName string) ([]string, error) {
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

func GetServiceVersionsWithHash(path, appName string) (map[string]string, error) {
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
		name, version = parts[len(parts)-2], parts[len(parts)-1]
	}

	return
}

func GetServiceVersionsFromCurrentBranch(path, appName string) ([]string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get all versions from repository
	allVersions, err := GetServiceVersionsWithHash(path, appName)
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
