package clients

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func UnitTests(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/clients/go")).
		// On package
		WithExec([]string{"go", "test", "./..."})
}
