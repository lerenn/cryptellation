// A generated module for CryptellationCandlesticksCi functions
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
	"cryptellation/svc/candlesticks/build/ci/dagger/internal/dagger"
)

const (
	path = "svc/candlesticks"
)

type CryptellationCandlesticksCi struct{}

func (mod *CryptellationCandlesticksCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().Linter(sourceDir, path)
}

func (mod *CryptellationCandlesticksCi) CheckGeneration(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().CheckGeneration(sourceDir, path)
}

func (mod *CryptellationCandlesticksCi) UnitTests(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationInternal().
		CryptellationGoCodeContainer(sourceDir, path).
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e ./internal/adapters -e ./test)",
		})
}

func (mod *CryptellationCandlesticksCi) IntegrationTests(
	ctx context.Context,
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	c := dag.CryptellationInternal().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationInternal().AttachMongo(c, dag.CryptellationInternal().Mongo().AsService())
	c = dag.CryptellationInternal().AttachBinance(c, secretsFile)
	return c.WithExec([]string{"sh", "-c",
		"go test ./internal/adapters/...",
	})
}

func (mod *CryptellationCandlesticksCi) EndToEndTests(
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	// Dependencies
	mongo := dag.CryptellationInternal().Mongo().AsService()
	nats := dag.CryptellationInternal().Nats().AsService()

	// Service
	candlesticks := dag.CryptellationCandlesticks().RunnerWithDependencies(
		sourceDir,
		secretsFile,
		mongo,
		nats,
	).AsService()

	// // Tester
	c := dag.CryptellationInternal().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationInternal().AttachNats(c, nats)
	c = c.WithServiceBinding("cryptellation-candlesticks", candlesticks)
	return c.WithExec([]string{
		"go", "test", "./test/...",
	})
}
