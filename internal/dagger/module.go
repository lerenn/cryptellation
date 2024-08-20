// A generated module for CryptellationInternal functions
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

	"github.com/lerenn/cryptellation/internal/dagger/internal/dagger"
)

type CryptellationInternal struct{}

// Linter returns a container that runs the linter.
func (mod *CryptellationInternal) Linter(sourceDir *dagger.Directory, path string) *dagger.Container {
	c := dag.Container().
		From("golangci/golangci-lint:v1.59.1").
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint"))

	c = mod.WithGoCodeAndCacheAsWorkDirectory(c, sourceDir, path)

	return c.WithExec([]string{"golangci-lint", "run", "--timeout", "10m"})
}

func (mod *CryptellationInternal) CheckGeneration(
	ctx context.Context,
	sourceDir *dagger.Directory,
	path string,
) *dagger.Container {
	return mod.CryptellationGoCodeContainer(sourceDir, path).
		WithExec([]string{"sh", "/go/src/github.com/lerenn/cryptellation/scripts/check-generation.sh"})
}

func (mod *CryptellationInternal) UnitTests(sourceDir *dagger.Directory, path string) *dagger.Container {
	return mod.CryptellationGoCodeContainer(sourceDir, path).
		WithExec([]string{"go", "test", "./..."})
}
