package ci

import (
	"dagger.io/dagger"
)

// Linter returns a container that runs the linter.
func Linter(client *dagger.Client, path string) *dagger.Container {
	return client.Container().
		From("golangci/golangci-lint:v1.55.2").
		With(SourceAsWorkdir(client, path)).
		WithMountedCache("/root/.cache/golangci-lint", client.CacheVolume("golangci-lint")).
		WithExec([]string{"golangci-lint", "run"})
}
