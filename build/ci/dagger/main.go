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

func (m *CryptellationCi) Lint(sourceDir *dagger.Directory) []*dagger.Container {
	return []*dagger.Container{
		// Client
		dag.CryptellationPkg().Linter(sourceDir, "./clients/go"),

		// Commands
		dag.CryptellationPkg().Linter(sourceDir, "./cmd/cryptellation"),
		dag.CryptellationPkg().Linter(sourceDir, "./cmd/cryptellation-tui"),

		// Examples
		dag.CryptellationPkg().Linter(sourceDir, "./examples/go"),

		// Package
		dag.CryptellationPkg().Linter(sourceDir, "./pkg"),

		// Services
		dag.CryptellationCandlesticksCi().Linter(sourceDir),
		dag.CryptellationBacktestsCi().Linter(sourceDir),
		dag.CryptellationExchangesCi().Linter(sourceDir),
		dag.CryptellationForwardtestsCi().Linter(sourceDir),
		dag.CryptellationIndicatorsCi().Linter(sourceDir),
		dag.CryptellationTicksCi().Linter(sourceDir),
	}
}
