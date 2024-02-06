package git

import "github.com/go-git/go-git/v5"

func ActualBranchName(path string) (string, error) {
	// Open git repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Name().Short(), nil
}
