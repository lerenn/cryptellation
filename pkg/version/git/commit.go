package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func lastCommit(path string) (*object.Commit, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get commit iterator
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	commit := &object.Commit{}
	if err := commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.After(commit.Author.When) {
			commit = c
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return commit, nil
}

// LastCommitHash returns the hash of the last commit in the git repository at path.
func LastCommitHash(path string) (string, error) {
	commit, err := lastCommit(path)
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

// LastCommitMessage returns the message of the last commit in the git repository at path.
func LastCommitMessage(path string) (string, error) {
	commit, err := lastCommit(path)
	if err != nil {
		return "", err
	}

	return commit.Message, nil
}
