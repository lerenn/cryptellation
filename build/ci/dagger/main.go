// A generated module for CryptellationCi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"

	"github.com/lerenn/cryptellation/v1/build/ci/dagger/internal/dagger"
)

type CryptellationCi struct{}

// Linter returns a container that runs the linter.
func (mod *CryptellationCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	c := dag.Container().
		From("golangci/golangci-lint:v1.62.0").
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint"))

	c = mod.withGoCodeAndCacheAsWorkDirectory(c, sourceDir)

	return c.WithExec([]string{"golangci-lint", "run", "--timeout", "10m"})
}

func (mod *CryptellationCi) CheckTODOs(
	sourceDir *dagger.Directory,
) *dagger.Container {
	return mod.cryptellationGoCodeContainer(sourceDir).
		WithExec([]string{"go", "run", "./tools/invtodos"})
}

func (mod *CryptellationCi) CheckGeneration(
	ctx context.Context,
	sourceDir *dagger.Directory,
) []*dagger.Container {
	return []*dagger.Container{
		mod.cryptellationGoCodeContainer(sourceDir).
			WithExec([]string{"sh", "-c",
				"go generate ./... && " +
					"sh scripts/check-generation.sh"}),
		mod.cryptellationPythonCodeContainer(sourceDir).
			WithExec([]string{"sh", "-c",
				"pip install -Ur clients/python/gateway/requirements.dev.txt && " +
					"sh scripts/check-generation.sh"}),
	}
}

func (mod *CryptellationCi) UnitTests(sourceDir *dagger.Directory) *dagger.Container {
	return mod.cryptellationGoCodeContainer(sourceDir).
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e /activities -e /test)",
		})
}

// Create a new release
func (ci *CryptellationCi) CreateRelease(
	ctx context.Context,
	sourceDir *dagger.Directory,
	// +optional
	sshPrivateKeyFile *dagger.Secret,
	// +optional
	cryptellationGitToken *dagger.Secret,
	// +optional
	cryptellationPullRequestToken *dagger.Secret,
) error {
	repo, err := NewGit(ctx, sourceDir, authParams{
		SSHPrivateKeyFile:             sshPrivateKeyFile,
		CryptellationPullRequestToken: cryptellationPullRequestToken,
		CryptellationGitToken:         cryptellationGitToken,
	})
	if err != nil {
		return err
	}

	// Update Source Code
	sourceDir, err = updateSourceCode(ctx, sourceDir, &repo)
	if err != nil {
		return err
	}
	repo.UpdateSourceDir(sourceDir)

	// Update Helm chart
	sourceDir, err = updateHelmChart(ctx, sourceDir, &repo)
	if err != nil {
		return err
	}
	repo.UpdateSourceDir(sourceDir)

	// Push new commit with tag
	return repo.PublishNewVersionAsCommit(ctx)
}

// Publish a new release
func (ci *CryptellationCi) PublishRelease(
	ctx context.Context,
	sourceDir *dagger.Directory,
	// +optional
	sshPrivateKeyFile *dagger.Secret,
	// +optional
	cryptellationGitToken *dagger.Secret,
	// +optional
	packagesGitToken *dagger.Secret,
) error {
	// Create auth params
	auth := authParams{
		SSHPrivateKeyFile:     sshPrivateKeyFile,
		PackagesGitToken:      packagesGitToken,
		CryptellationGitToken: cryptellationGitToken,
	}

	// Create Git repo access
	repo, err := NewGit(ctx, sourceDir, auth)
	if err != nil {
		return err
	}

	// Publish new tag
	if err := repo.PublishTagFromReleaseTitle(ctx); err != nil {
		return err
	}

	// Publish Docker images
	if err := publishDockerImages(ctx, sourceDir, &repo); err != nil {
		return err
	}

	// Publish Helm chart
	return publishHelmChart(ctx, sourceDir, &repo, auth)
}
