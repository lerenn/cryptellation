package main

import (
	"runtime"

	"github.com/lerenn/cryptellation/v1/build/ci/dagger/internal/dagger"
)

func (mod *CryptellationCi) cryptellationPythonCodeContainer(
	sourceDir *dagger.Directory,
) *dagger.Container {
	c := dag.Container().From("python:3-alpine")
	return mod.withPythonCodeAndCacheAsWorkDirectory(c, sourceDir)
}

func (mod *CryptellationCi) withPythonCodeAndCacheAsWorkDirectory(
	c *dagger.Container,
	sourceDir *dagger.Directory,
) *dagger.Container {
	return c.
		// Add Go caches
		WithMountedCache("/root/.cache/pip", dag.CacheVolume("pip-cache")).

		// Add source code
		WithMountedDirectory("/src/github.com/lerenn/cryptellation", sourceDir).

		// Add workdir
		WithWorkdir("/src/github.com/lerenn/cryptellation")
}

func (mod *CryptellationCi) cryptellationGoCodeContainer(
	sourceDir *dagger.Directory,
) *dagger.Container {
	c := dag.Container().From("golang:" + goVersion() + "-alpine")
	return mod.withGoCodeAndCacheAsWorkDirectory(c, sourceDir)
}

func (mod *CryptellationCi) withGoCodeAndCacheAsWorkDirectory(
	c *dagger.Container,
	sourceDir *dagger.Directory,
) *dagger.Container {
	return c.
		// Add Go caches
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild")).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gocache")).

		// Add source code
		WithMountedDirectory("/go/src/github.com/lerenn/cryptellation", sourceDir).

		// Add workdir
		WithWorkdir("/go/src/github.com/lerenn/cryptellation")
}

func goVersion() string {
	return runtime.Version()[2:]
}
