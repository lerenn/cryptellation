package git

import "github.com/go-git/go-git/v5"

// ActualBranchName returns the actual branch name of the git repository at path.
func ActualBranchName(path string) (string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Name().Short(), nil
}
