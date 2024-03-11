package ci

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client, path string) func(ctx context.Context) error {
	if path[0] == '/' {
		path = "." + path
	}

	return func(ctx context.Context) error {
		_, err := client.Container().
			// Add base image
			From("golang:"+utils.GoVersion()+"-alpine3.19").
			// Add source code as work directory
			With(SourceAsWorkdir(client, path)).
			// Add command to generate code
			WithExec([]string{"go", "generate", "./..."}).
			// Export directory
			Directory(".").
			Export(ctx, path)

		return err
	}
}
