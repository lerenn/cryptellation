package main

import (
	"context"
	"dagger/cryptellation-ci/internal/dagger"
	"strings"
)

type Git struct {
	container *dagger.Container

	lastCommit struct {
		title    string
		shortSHA string
	}

	lastTag      string
	newSemVer    string
	actualBranch string
}

func NewGit(srcDir *dagger.Directory, sshPrivateKeyFile *dagger.Secret) Git {
	container := dag.Container().
		From("alpine/git").
		WithMountedDirectory("/git", srcDir).
		WithWorkdir("/git").
		WithoutEntrypoint()

	if sshPrivateKeyFile != nil {
		container = container.WithMountedSecret("/root/.ssh/id_rsa", sshPrivateKeyFile)
	}

	return Git{
		container: container,
	}
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

func (g *Git) PushNewCommitWithTag(ctx context.Context) error {
	// Get new semver
	semver, err := g.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return err
	}

	// Check if semver is empty
	if semver == "" {
		return nil
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
		WithExec([]string{"git", "commit", "-m", "chore: release " + semver}).
		Sync(ctx)
	if err != nil {
		return err
	}
	g.resetLastCommitCache()

	// Push new commit
	g.container, err = g.container.
		WithExec([]string{"git", "push"}).
		Sync(ctx)
	if err != nil {
		return err
	}

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
