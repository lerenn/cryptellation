// A generated module for CryptellationTicksCi functions
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
	"cryptellation/svc/ticks/build/ci/dagger/internal/dagger"
)

const (
	path = "svc/ticks"
)

type CryptellationTicksCi struct{}

func (mod *CryptellationTicksCi) Linter(sourceDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().Linter(sourceDir, path)
}

func (mod *CryptellationTicksCi) CheckGeneration(rootDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().CheckGeneration(rootDir, path)
}

func (mod *CryptellationTicksCi) UnitTests(rootDir *dagger.Directory) *dagger.Container {
	return dag.CryptellationPkg().
		CryptellationGoCodeContainer(rootDir, path).
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e ./internal/adapters -e ./test)",
		})
}

func (mod *CryptellationTicksCi) IntegrationTests(
	ctx context.Context,
	rootDir *dagger.Directory,
	secretsFile *dagger.Secret,
) *dagger.Container {
	c := dag.CryptellationPkg().CryptellationGoCodeContainer(rootDir, path)
	c = dag.CryptellationPkg().AttachMongo(c, dag.CryptellationPkg().Mongo().AsService())
	c = dag.CryptellationPkg().AttachNats(c, dag.CryptellationPkg().Nats().AsService())
	c = dag.CryptellationPkg().AttachBinance(c, secretsFile)
	return c.WithExec([]string{"sh", "-c",
		"go test ./internal/adapters/...",
	})
}
