package ci

import (
	"context"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
	"github.com/joho/godotenv"
)

const (
	secretsFilePath = "./.credentials.env"
)

func SourceAsWorkdir(client *dagger.Client, path string) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add Go caches
			WithMountedCache("/root/.cache/go-build", client.CacheVolume("gobuild")).
			WithMountedCache("/go/pkg/mod", client.CacheVolume("gocache")).

			// Add source code
			WithMountedDirectory("/go/src/github.com/lerenn/cryptellation", client.Host().Directory(".")).

			// Add workdir
			WithWorkdir("/go/src/github.com/lerenn/cryptellation" + path)
	}
}

func Secret(client *dagger.Client, name string) func(r *dagger.Container) *dagger.Container {
	// Check if secrets file exists
	if _, err := os.Stat(secretsFilePath); err != nil {
		panic(err)
	}

	// Load file with secrets
	envMap, err := godotenv.Read(secretsFilePath)
	if err != nil {
		panic(err)
	}

	// Get requested secret from loaded file
	content, exists := envMap[name]
	if !exists {
		panic(fmt.Errorf("there is no %q variable in %q", name, secretsFilePath))
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
	for _, l1 := range containers {
		for _, l2 := range l1 {
			if l2 == nil {
				continue
			}

			rContainers = append(rContainers, l2)
		}
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
