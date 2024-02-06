package git

import "github.com/go-git/go-git/v5"

func Tag(path, tag string) error {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	// Get last commit
	commit, err := lastCommit(path)
	if err != nil {
		return err
	}

	// Create tag
	_, err = repo.CreateTag(tag, commit.Hash, nil)
	return err
}
