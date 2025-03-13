package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/lerenn/cryptellation/v1/build/ci/dagger/internal/dagger"

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
	actualBranch string

	semVerChange  SemVerChange
	updatedSemVer string
}

type authParams struct {
	// SSH Private Key File mode (ie. access to everything)
	SSHPrivateKeyFile *dagger.Secret
	// Token mode (ie. ine grained access)
	CryptellationGitToken         *dagger.Secret
	CryptellationPullRequestToken *dagger.Secret
	PackagesGitToken              *dagger.Secret
}

func NewGit(ctx context.Context, srcDir *dagger.Directory, auth authParams) (Git, error) {
	var err error

	// Create container
	container := dag.Container().
		From("alpine/git").
		WithMountedDirectory("/git", srcDir).
		WithWorkdir("/git").
		WithoutEntrypoint()

	// Set authentication based on the provided parameters
	if auth.SSHPrivateKeyFile != nil {
		container, err = mountSSHPrivateKeyFile(ctx, container, auth)
		if err != nil {
			return Git{}, err
		}
	} else if auth.CryptellationGitToken != nil {
		token, err := auth.CryptellationGitToken.Plaintext(ctx)
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
	} else {
		return Git{}, fmt.Errorf("no auth method provided")
	}

	// Set Git author
	container, err = setGitAuthor(ctx, container)
	if err != nil {
		return Git{}, err
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

func (g *Git) GetNewSemVerIfNeeded(ctx context.Context) (change SemVerChange, semVer string, err error) {
	// Check if already doesn't exist
	if g.semVerChange != SemVerChangeUnknown {
		return g.semVerChange, g.updatedSemVer, nil
	}

	// Get last commit title
	title, err := g.GetLastCommitTitle(ctx)
	if err != nil {
		return SemVerChangeUnknown, "", err
	}

	// Get last tag
	tag, err := g.GetLastTag(ctx)
	if err != nil {
		return SemVerChangeUnknown, "", err
	}

	// Generate new version based on title and return if there is no change
	change, semVer, err = processSemVerChange(tag, title)
	if err != nil {
		return SemVerChangeUnknown, "", err
	}

	// Add to cache
	g.semVerChange = change
	g.updatedSemVer = semVer

	return change, semVer, nil
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
	title = strings.TrimSuffix(title, "\n")
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
	commitMsg := fmt.Sprintf("%q", title)
	g.container, err = g.container.
		WithExec([]string{"git", "commit", "-m", commitMsg}).
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
	if g.auth.CryptellationPullRequestToken != nil {
		token, err := g.auth.CryptellationPullRequestToken.Plaintext(ctx)
		if err != nil {
			return err
		}

		client := github.NewClient(nil).WithAuthToken(token)
		if _, _, err := client.PullRequests.Create(ctx, "lerenn", "cryptellation", &github.NewPullRequest{
			Title:               &title,
			Base:                toReference("main"),
			Head:                toReference(branchName),
			MaintainerCanModify: toReference(true),
		}); err != nil {
			return err
		}
	}

	return nil
}

func mountSSHPrivateKeyFile(
	ctx context.Context,
	container *dagger.Container,
	auth authParams,
) (*dagger.Container, error) {
	// Mount key
	if auth.SSHPrivateKeyFile != nil {
		container = container.WithMountedSecret("/root/.ssh/id_rsa", auth.SSHPrivateKeyFile)
	}

	// Add SSH config for github
	container, err := container.
		WithExec([]string{"sh", "-c", "echo -e 'Host github.com\n\tStrictHostKeyChecking no\n' > /root/.ssh/config"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func setGitAuthor(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	// Add infos on author
	container, err := container.
		WithExec([]string{"git", "config", "--global", "user.email", "louis.fradin+cryptellation-ci@gmail.com"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}
	container, err = container.
		WithExec([]string{"git", "config", "--global", "user.name", "Cryptellation CI"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (g *Git) PublishNewVersionAsCommit(ctx context.Context) error {
	// Get new semver
	change, newSemVer, err := g.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return err
	} else if change == SemVerChangeNone {
		return nil
	}

	// Push new commit with tag
	return g.PublishNewCommit(ctx, "release: v"+newSemVer)
}
