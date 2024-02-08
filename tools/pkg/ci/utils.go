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
			WithMountedDirectory(
				"/go/src/github.com/lerenn/cryptellation",
				client.Host().Directory(".", dagger.HostDirectoryOpts{Exclude: []string{"tmp", "website"}})).

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
	funcs := make([]func(context.Context) error, 0)
	for _, l1 := range containers {
		for _, l2 := range l1 {
			if l2 == nil {
				continue
			}

			// Note: create a new local variable to store value of actual l2
			callback := l2

			fn := func(ctx context.Context) error {
				_, err := callback.Stderr(ctx)
				return err
			}
			funcs = append(funcs, fn)
		}
	}

	ExecuteInParallel(ctx, funcs...)
}

func ExecuteInParallel(ctx context.Context, funcs ...func(context.Context) error) {
	// Excute containers
	var wg sync.WaitGroup
	for _, fn := range funcs {
		go func(callback func(context.Context) error) {
			if err := callback(ctx); err != nil {
				panic(err)
			}
			wg.Done()
		}(fn)

		wg.Add(1)
	}

	wg.Wait()
}

func ExposeOnLocalPort(
	client *dagger.Client,
	service *dagger.Service,
	ports ...dagger.PortForward,
) func(ctx context.Context, opts ...dagger.ServiceStopOpts) (*dagger.Service, error) {
	// Set tunnel to service
	tunnel, err := client.Host().Tunnel(service, dagger.HostTunnelOpts{
		Ports: ports,
	}).Start(context.Background())
	if err != nil {
		panic(err)
	}

	return tunnel.Stop
}
