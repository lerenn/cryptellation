package ci

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
	"github.com/lerenn/cryptellation/tools/pkg/ci"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, err := client.Container().
			// Add base image
			From("golang:"+utils.GoVersion()+"-alpine3.19").
			// Add source code as work directory
			With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
			// Add command to generate code
			WithExec([]string{"go", "generate", "./..."}).
			// Export directory
			Directory(".").
			Export(ctx, "./svc/"+ServiceName)

		return err
	}
}
