package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/utils"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Add command to generate code
		WithExec([]string{"go", "generate", "./..."})
}
