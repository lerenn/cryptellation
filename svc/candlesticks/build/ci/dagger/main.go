// A generated module for CryptellationCandlesticksCi functions
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

import "cryptellation/svc/candlesticks/build/ci/dagger/internal/dagger"

const (
	path = "svc/candlesticks"
)

type CryptellationCandlesticksCi struct{}

func (mod *CryptellationCandlesticksCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().Linter(sourceDir, path)
}

func (mod *CryptellationCandlesticksCi) CheckGeneration(rootDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().CheckGeneration(rootDir, path)
}
