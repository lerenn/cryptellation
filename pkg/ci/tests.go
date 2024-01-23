package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func UnitTests(client *dagger.Client, path string) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(SourceAsWorkdir(client, path)).
		// On package
		WithExec([]string{"go", "test", "./..."})
}
