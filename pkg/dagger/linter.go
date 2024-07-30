package main

import "cryptellation/pkg/dagger/internal/dagger"

// Linter returns a container that runs the linter.
func (mod *CryptellationPkg) Linter(sourceDir *dagger.Directory, path string) *dagger.Container {
	c := dag.Container().
		From("golangci/golangci-lint:v1.59.1").
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint"))

	c = mod.WithGoCodeAndCacheAsWorkDirectory(c, sourceDir, path)

	return c.WithExec([]string{"golangci-lint", "run", "--timeout", "10m"})
}
