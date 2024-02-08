package ci

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
)

// UpdateGoMod returns a container that generates code.
func UpdateGoMod(client *dagger.Client, path string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, err := client.Container().
			// Add base image
			From("golang:"+utils.GoVersion()+"-alpine3.19").
			// Add source code as work directory
			With(SourceAsWorkdir(client, "/"+path)).
			// Add command to generate code
			WithExec([]string{"go", "mod", "tidy"}).
			// Export directory
			Directory(".").
			Export(ctx, path)

		return err
	}
}
