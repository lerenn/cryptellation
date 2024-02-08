package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func UnitTests(client *dagger.Client, path string) *dagger.Container {
	return client.Container().
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		With(SourceAsWorkdir(client, path)).
		WithExec([]string{"go", "test", "./..."})
}
