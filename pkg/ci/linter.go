package ci

import (
	"dagger.io/dagger"
)

// Linter returns a container that runs the linter.
func Linter(client *dagger.Client, path string) *dagger.Container {
	return client.Container().
		// Add base image
		From("golangci/golangci-lint:v1.55.2").
		// Add source code as work directory
		With(SourceAsWorkdir(client, path)).
		// Add golangci-lint cache
		WithMountedCache("/root/.cache/golangci-lint", client.CacheVolume("golangci-lint")).
		// Add command
		WithExec([]string{"golangci-lint", "run"})
}
