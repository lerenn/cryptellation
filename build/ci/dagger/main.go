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

	"github.com/lerenn/cryptellation/build/ci/dagger/internal/dagger"
)

type CryptellationCi struct{}

// Check returns a list of containers for checking everything
func (m *CryptellationCi) Check(sourceDir *dagger.Directory, secretsFile *dagger.Secret) []*dagger.Container {
	containers := make([]*dagger.Container, 0)
	containers = append(containers, m.Lint(sourceDir)...)
	containers = append(containers, m.CheckGeneration(sourceDir)...)
	containers = append(containers, m.UnitTests(sourceDir)...)
	containers = append(containers, m.IntegrationTests(sourceDir, secretsFile)...)
	containers = append(containers, m.EndToEndTests(sourceDir, secretsFile)...)
	return containers
}

// Lint returns a list of containers for linting
func (m *CryptellationCi) Lint(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().Linter(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().Linter(sourceDir, "./cmd/cryptellation"),

		// Examples
		dag.CryptellationInternal().Linter(sourceDir, "./examples/go"),

		// Internal
		dag.CryptellationInternal().Linter(sourceDir, "./internal"),

		// Package
		dag.CryptellationInternal().Linter(sourceDir, "./pkg"),

		// Services
		dag.CryptellationCandlesticksCi().Linter(sourceDir),
		dag.CryptellationBacktestsCi().Linter(sourceDir),
		dag.CryptellationExchangesCi().Linter(sourceDir),
		dag.CryptellationForwardtestsCi().Linter(sourceDir),
		dag.CryptellationIndicatorsCi().Linter(sourceDir),
		dag.CryptellationTicksCi().Linter(sourceDir),
	}
}

// CheckGeneration returns a list of containers for checking generation
func (m *CryptellationCi) CheckGeneration(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./cmd/cryptellation"),

		// Examples
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./examples/go"),

		// Internal
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./internal"),

		// Package
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./pkg"),

		// Services
		dag.CryptellationCandlesticksCi().CheckGeneration(sourceDir),
		dag.CryptellationBacktestsCi().CheckGeneration(sourceDir),
		dag.CryptellationExchangesCi().CheckGeneration(sourceDir),
		dag.CryptellationForwardtestsCi().CheckGeneration(sourceDir),
		dag.CryptellationIndicatorsCi().CheckGeneration(sourceDir),
		dag.CryptellationTicksCi().CheckGeneration(sourceDir),
	}
}

// UnitTests returns a list of containers for unit tests
func (m *CryptellationCi) UnitTests(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().UnitTests(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().UnitTests(sourceDir, "./cmd/cryptellation"),

		// Examples
		dag.CryptellationInternal().UnitTests(sourceDir, "./examples/go"),

		// Internal
		dag.CryptellationInternal().UnitTests(sourceDir, "./internal"),

		// Package
		dag.CryptellationInternal().UnitTests(sourceDir, "./pkg"),

		// Services
		dag.CryptellationCandlesticksCi().UnitTests(sourceDir),
		dag.CryptellationBacktestsCi().UnitTests(sourceDir),
		dag.CryptellationExchangesCi().UnitTests(sourceDir),
		dag.CryptellationForwardtestsCi().UnitTests(sourceDir),
		dag.CryptellationIndicatorsCi().UnitTests(sourceDir),
		dag.CryptellationTicksCi().UnitTests(sourceDir),
	}
}

// IntegrationTests returns a list of containers for integration tests
func (m *CryptellationCi) IntegrationTests(sourceDir *dagger.Directory, secretsFile *dagger.Secret) []*dagger.Container {
	return []*dagger.Container{
		dag.CryptellationBacktestsCi().IntegrationTests(sourceDir),
		dag.CryptellationCandlesticksCi().IntegrationTests(sourceDir, secretsFile),
		dag.CryptellationExchangesCi().IntegrationTests(sourceDir, secretsFile),
		dag.CryptellationForwardtestsCi().IntegrationTests(sourceDir),
		dag.CryptellationIndicatorsCi().IntegrationTests(sourceDir),
		dag.CryptellationTicksCi().IntegrationTests(sourceDir, secretsFile),
	}
}

// EndToEndTests returns a list of containers for end-to-end tests
func (m *CryptellationCi) EndToEndTests(sourceDir *dagger.Directory, secretsFile *dagger.Secret) []*dagger.Container {
	return []*dagger.Container{
		dag.CryptellationBacktestsCi().EndToEndTests(sourceDir, secretsFile),
		dag.CryptellationCandlesticksCi().EndToEndTests(sourceDir, secretsFile),
		dag.CryptellationExchangesCi().EndToEndTests(sourceDir, secretsFile),
		dag.CryptellationForwardtestsCi().EndToEndTests(sourceDir, secretsFile),
		dag.CryptellationIndicatorsCi().EndToEndTests(sourceDir, secretsFile),
		dag.CryptellationTicksCi().EndToEndTests(sourceDir, secretsFile),
	}
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

	// Update Helm chart
	sourceDir, err = updateHelmChart(ctx, sourceDir, &repo)
	if err != nil {
		return err
	}
	repo.UpdateSourceDir(sourceDir)

	// Get new semver
	newSemVer, err := repo.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return err
	} else if newSemVer == "" {
		return nil
	}

	// Push new commit with tag
	return repo.PublishNewCommit(ctx, "release: "+newSemVer)
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
