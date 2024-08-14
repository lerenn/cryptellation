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
	"dagger/cryptellation-ci/internal/dagger"
)

type CryptellationCi struct{}

func (m *CryptellationCi) Check(sourceDir *dagger.Directory, secretsFile *dagger.Secret) []*dagger.Container {
	containers := make([]*dagger.Container, 0)
	containers = append(containers, m.Lint(sourceDir)...)
	containers = append(containers, m.CheckGeneration(sourceDir)...)
	containers = append(containers, m.UnitTests(sourceDir)...)
	containers = append(containers, m.IntegrationTests(sourceDir, secretsFile)...)
	containers = append(containers, m.EndToEndTests(sourceDir, secretsFile)...)
	return containers
}

func (m *CryptellationCi) Lint(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().Linter(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().Linter(sourceDir, "./cmd/cryptellation"),
		dag.CryptellationInternal().Linter(sourceDir, "./cmd/cryptellation-tui"),

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

func (m *CryptellationCi) CheckGeneration(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./cmd/cryptellation"),
		dag.CryptellationInternal().CheckGeneration(sourceDir, "./cmd/cryptellation-tui"),

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

func (m *CryptellationCi) UnitTests(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationInternal().UnitTests(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationInternal().UnitTests(sourceDir, "./cmd/cryptellation"),
		dag.CryptellationInternal().UnitTests(sourceDir, "./cmd/cryptellation-tui"),

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
