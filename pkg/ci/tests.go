package ci

import (
	"cryptellation/pkg/utils"

	"dagger.io/dagger"
)

func UnitTests(client *dagger.Client, path string) *dagger.Container {
	return client.Container().
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		With(SourceAsWorkdir(client, path)).
		WithExec([]string{"go", "test", "./..."})
}
