package git

import "github.com/go-git/go-git/v5"

func LastCommitHash(path string) (string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	// Get commit iterator
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return "", err
	}

	// Get next commit
	commit, err := commitIter.Next()
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}
