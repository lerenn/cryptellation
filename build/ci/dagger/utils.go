package main

import (
	"context"
	"dagger/cryptellation-ci/internal/dagger"
)

func setGithubRepositoryRequirements(
	ctx context.Context,
	container *dagger.Container,
) (*dagger.Container, error) {
	// Add SSH config for github
	container, err := container.
		WithExec([]string{"sh", "-c", "echo -e 'Host github.com\n\tStrictHostKeyChecking no\n' > /root/.ssh/config"}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	// Add infos on author
	container, err = container.
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
