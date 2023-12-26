package pipeline

import (
	"context"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
	"github.com/joho/godotenv"
)

func sourceAsWorkdir(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add Go caches
			WithMountedCache("/root/.cache/go-build", client.CacheVolume("gobuild")).
			WithMountedCache("/go/pkg/mod", client.CacheVolume("gocache")).

			// Add source code
			WithMountedDirectory(SourcePath, client.Host().Directory(".")).

			// Add workdir
			WithWorkdir(SourcePath)
	}
}

func containerWithSourceCodeAndGo(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From(GolangImage).
		// Add source code as work directory
		With(sourceAsWorkdir(client))
}

func secret(client *dagger.Client, name string) func(r *dagger.Container) *dagger.Container {
	// Check if secrets file exists
	if _, err := os.Stat(SecretsFilePath); err != nil {
		panic(err)
	}

	// Load file with secrets
	envMap, err := godotenv.Read(SecretsFilePath)
	if err != nil {
		panic(err)
	}

	// Get requested secret from loaded file
	content, exists := envMap[name]
	if !exists {
		panic(fmt.Errorf("there is no %q variable in %q", name, SecretsFilePath))
	}

	// Change to dagger secret
	daggerSecret := client.SetSecret(name, content)

	// Add it to the container
	return func(r *dagger.Container) *dagger.Container {
		return r.WithSecretVariable(name, daggerSecret)
	}
}

func ExecuteContainersInParallel(ctx context.Context, containers ...[]*dagger.Container) {
	// Regroup arg
	rContainers := make([]*dagger.Container, 0)
	for _, c := range containers {
		rContainers = append(rContainers, c...)
	}

	// Excute containers
	var wg sync.WaitGroup
	for _, ec := range rContainers {
		go func(e *dagger.Container) {
			_, err := e.Stderr(ctx)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(ec)

		wg.Add(1)
	}

	wg.Wait()
}
