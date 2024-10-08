// A generated module for CryptellationForwardtestsCi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/svc/forwardtests/build/ci/dagger/internal/dagger"

	"github.com/lerenn/cryptellation/internal/docker"
)

const (
	path            = "svc/forwardtests"
	dockerImageName = "lerenn/cryptellation-forwardtests"
)

type CryptellationForwardtestsCi struct{}

func (mod *CryptellationForwardtestsCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().Linter(sourceDir, path)
}

func (mod *CryptellationForwardtestsCi) CheckGeneration(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().CheckGeneration(sourceDir, path)
}

func (mod *CryptellationForwardtestsCi) UnitTests(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().
		CryptellationGoCodeContainer(sourceDir, path).
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e ./internal/adapters -e ./test)",
		})
}

func (mod *CryptellationForwardtestsCi) IntegrationTests(
	ctx context.Context,
	sourceDir *dagger.Directory,
) *dagger.Container {
	c := dag.CryptellationInternal().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationInternal().AttachMongo(c, dag.CryptellationInternal().Mongo().AsService())
	c = dag.CryptellationInternal().AttachNats(c, dag.CryptellationInternal().Nats().AsService())
	return c.WithExec([]string{"sh", "-c",
		"go test ./internal/adapters/...",
	})
}

func (mod *CryptellationForwardtestsCi) EndToEndTests(
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	// Dependencies
	mongo := dag.CryptellationInternal().Mongo().AsService()
	nats := dag.CryptellationInternal().Nats().AsService()
	candlesticks := dag.CryptellationCandlesticks().
		RunnerWithDependencies(sourceDir, secretsFile, mongo, nats).
		AsService()

	// Service
	forwardtestsService := dag.CryptellationForwardtests().
		RunnerWithDependencies(sourceDir, candlesticks, mongo, nats).
		AsService()

	// Tester
	c := dag.CryptellationInternal().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationInternal().AttachNats(c, nats)
	c = c.WithServiceBinding("cryptellation-forwardtests", forwardtestsService)
	return c.WithExec([]string{
		"go", "test", "./test/...",
	})
}

// Publishes the Docker image
func (ci *CryptellationForwardtestsCi) PublishDockerImage(
	ctx context.Context,
	sourceDir *dagger.Directory,
	tags []string,
) error {
	// Get images for each platform
	platformVariants := make([]*dagger.Container, 0, len(docker.GoRunnersInfo))
	for targetPlatform := range docker.GoRunnersInfo {
		runner := dag.CryptellationForwardtests().Runner(sourceDir, dagger.CryptellationForwardtestsRunnerOpts{
			TargetPlatform: targetPlatform,
		})

		platformVariants = append(platformVariants, runner)
	}

	// Set publication options from images
	publishOpts := dagger.ContainerPublishOpts{
		PlatformVariants: platformVariants,
	}

	// Publish with tags
	for _, tag := range tags {
		addr := fmt.Sprintf("%s:%s", dockerImageName, tag)
		if _, err := dag.Container().Publish(ctx, addr, publishOpts); err != nil {
			return err
		}
	}

	return nil
}
