package pipeline

import (
	"dagger.io/dagger"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client) *dagger.Container {
	return containerWithSourceCodeAndGo(client).
		// Add command to generate code
		WithExec([]string{"go", "generate", "./..."})
}
