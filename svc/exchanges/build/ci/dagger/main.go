// A generated module for CryptellationExchangesCi functions
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
	"cryptellation/svc/exchanges/build/ci/dagger/internal/dagger"
)

const (
	path = "svc/exchanges"
)

type CryptellationExchangesCi struct{}

func (mod *CryptellationExchangesCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().Linter(sourceDir, path)
}

func (mod *CryptellationExchangesCi) CheckGeneration(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().CheckGeneration(sourceDir, path)
}

func (mod *CryptellationExchangesCi) UnitTests(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().
		CryptellationGoCodeContainer(sourceDir, path).
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e ./internal/adapters -e ./test)",
		})
}

func (mod *CryptellationExchangesCi) IntegrationTests(
	ctx context.Context,
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	c := dag.CryptellationPkg().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationPkg().AttachMongo(c, dag.CryptellationPkg().Mongo().AsService())
	c = dag.CryptellationPkg().AttachBinance(c, secretsFile)
	return c.WithExec([]string{"sh", "-c",
		"go test ./internal/adapters/...",
	})
}

func (mod *CryptellationExchangesCi) EndToEndTests(
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	// Dependencies
	mongo := dag.CryptellationPkg().Mongo().AsService()
	nats := dag.CryptellationPkg().Nats().AsService()

	// Service
	exchangesService := dag.CryptellationExchanges().RunnerWithDependencies(
		sourceDir,
		secretsFile,
		mongo,
		nats,
	).AsService()

	// Tester
	c := dag.CryptellationPkg().CryptellationGoCodeContainer(sourceDir, path)
	c = dag.CryptellationPkg().AttachNats(c, nats)
	c = c.WithServiceBinding("cryptellation-exchanges", exchangesService)
	return c.WithExec([]string{
		"go", "test", "./test/...",
	})
}
