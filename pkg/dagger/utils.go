package main

import (
	"cryptellation/pkg/dagger/internal/dagger"
)

func (mod *CryptellationPkg) WithGoCodeAndCacheAsWorkDirectory(
	c *dagger.Container,
	rootDir *dagger.Directory,
	path string,
) *dagger.Container {
	return c.
		// Add Go caches
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild")).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gocache")).

		// Add source code
		WithMountedDirectory("/go/src/cryptellation", rootDir).

		// Add workdir
		WithWorkdir("/go/src/cryptellation/" + path)
}

func (mod *CryptellationPkg) CryptellationGoCodeContainer(
	sourceDir *dagger.Directory,
	path string,
) *dagger.Container {
	c := dag.Container().From(golangImage)
	return mod.WithGoCodeAndCacheAsWorkDirectory(c, sourceDir, path)
}
