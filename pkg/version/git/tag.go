package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// TagCommit will create a tag on the last commit of the git repository at path.
func TagCommit(path, tag string) error {
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

// PushTags will execute a git push with tags.
func PushTags(path, tag string) error {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	return repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/tags/" + tag + ":refs/tags/" + tag),
		},
	})
}
