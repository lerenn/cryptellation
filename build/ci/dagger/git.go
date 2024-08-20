package main

import (
	"context"
	"cryptellation/pkg/utils"
	"dagger/cryptellation-ci/internal/dagger"
	"fmt"
	"strings"

	"github.com/google/go-github/v63/github"
)

type Git struct {
	container *dagger.Container
	auth      authParams

	lastCommit struct {
		title    string
		shortSHA string
	}

	lastTag      string
	newSemVer    string
	actualBranch string
}

type authParams struct {
	SSHPrivateKeyFile *dagger.Secret
	GithubToken       *dagger.Secret
}

func (ap authParams) Validate() error {
	if ap.SSHPrivateKeyFile == nil && ap.GithubToken == nil {
		return fmt.Errorf("either SSHPrivateKeyFile or GithubToken must be provided")
	}
	return nil
}

func NewGit(ctx context.Context, srcDir *dagger.Directory, auth authParams) (Git, error) {
	// Check auth params
	err := auth.Validate()
	if err != nil {
		return Git{}, err
	}

	// Create container
	container := dag.Container().
		From("alpine/git").
		WithMountedDirectory("/git", srcDir).
		WithWorkdir("/git").
		WithoutEntrypoint()

	// Set authentication based on the provided parameters
	if auth.SSHPrivateKeyFile != nil {
		// Mount the SSH private key
		container = container.WithMountedSecret("/root/.ssh/id_rsa", auth.SSHPrivateKeyFile)

		// Set github repository requirements
		container, err = setGithubRepositoryRequirements(ctx, container)
		if err != nil {
			return Git{}, err
		}
	} else if auth.GithubToken != nil {
		token, err := auth.GithubToken.Plaintext(ctx)
		if err != nil {
			return Git{}, err
		}

		// Change the url to use the token
		container, err = container.WithExec([]string{
			"git", "remote", "set-url", "origin", "https://lerenn:" + token + "@github.com/lerenn/cryptellation.git",
		}).Sync(ctx)
		if err != nil {
			return Git{}, err
		}
	}

	return Git{
		container: container,
		auth:      auth,
	}, nil
}

func (g *Git) UpdateSourceDir(srcDir *dagger.Directory) {
	g.container = g.container.WithoutMount("/git").WithMountedDirectory("/git", srcDir)
}

func (g *Git) GetLastCommitShortSHA(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastCommit.shortSHA != "" {
		return g.lastCommit.shortSHA, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "rev-parse", "--short", "HEAD"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	// Set the cache value
	g.lastCommit.shortSHA = res

	return g.lastCommit.shortSHA, nil
}

func (g *Git) GetActualBranch(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.actualBranch != "" {
		return g.actualBranch, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	// Set the cache value
	g.actualBranch = res

	return g.actualBranch, nil
}

func (g *Git) GetLastCommitTitle(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastCommit.title != "" {
		return g.lastCommit.title, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "log", "-1", "--pretty=%B"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	// Set the cache value
	g.lastCommit.title = res

	return g.lastCommit.title, nil
}

func (g *Git) GetLastTag(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastTag != "" {
		return g.lastTag, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Remove potential new line
	res = strings.TrimSuffix(res, "\n")

	// Set the cache value
	g.lastTag = res

	return g.lastTag, nil
}

func (g *Git) GetNewSemVerIfNeeded(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.newSemVer != "" {
		return g.newSemVer, nil
	}

	// Get last commit title
	title, err := g.GetLastCommitTitle(ctx)
	if err != nil {
		return "", err
	}

	// Get last tag
	tag, err := g.GetLastTag(ctx)
	if err != nil {
		return "", err
	}

	// Generate new version based on title and return if there is no change
	newSemVer, err := parseAndUpdateSemVer(tag, title)
	if err != nil {
		return "", err
	}
	newSemVer = "v" + newSemVer

	// Check if new version is the same as the last tag
	if newSemVer == tag {
		return "", nil
	}

	return newSemVer, nil
}

func (g *Git) resetLastCommitCache() {
	g.lastCommit.title = ""
	g.lastCommit.shortSHA = ""
}

func (g *Git) PublishTagFromReleaseTitle(ctx context.Context) error {
	// Get new semver
	title, err := g.GetLastCommitTitle(ctx)
	if err != nil {
		return err
	}
	title = strings.TrimPrefix(title, "release: ")
	semver := strings.Split(title, " ")[0]

	// Tag commit
	g.container, err = g.container.
		WithExec([]string{"git", "tag", semver}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Push new tag
	g.container, err = g.container.
		WithExec([]string{"git", "push", "--tags"}).
		Sync(ctx)

	return err
}

func (g *Git) PublishNewCommit(
	ctx context.Context,
	title string,
) error {
	var err error

	// Set new branch
	branchName := strings.ReplaceAll(title, " ", "-")
	branchName = strings.ReplaceAll(branchName, ":", "")
	g.container, err = g.container.
		WithExec([]string{"git", "checkout", "-b", branchName}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Add all changes
	g.container, err = g.container.
		WithExec([]string{"git", "add", "."}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Commit changes
	g.container, err = g.container.
		WithExec([]string{"git", "commit", "-m", title}).
		Sync(ctx)
	if err != nil {
		return err
	}
	g.resetLastCommitCache()

	// Push new commit
	g.container, err = g.container.
		WithExec([]string{"git", "push", "--set-upstream", "origin", branchName}).
		Sync(ctx)
	if err != nil {
		return err
	}

	// Create pull request
	if g.auth.GithubToken != nil {
		token, err := g.auth.GithubToken.Plaintext(ctx)
		if err != nil {
			return err
		}

		client := github.NewClient(nil).WithAuthToken(token)
		if _, _, err := client.PullRequests.Create(ctx, "lerenn", "cryptellation", &github.NewPullRequest{
			Title:               &title,
			Base:                utils.ToReference("main"),
			Head:                utils.ToReference(branchName),
			MaintainerCanModify: utils.ToReference(true),
		}); err != nil {
			return err
		}
	}

	return nil
}
