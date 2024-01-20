package pkg

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func UnitTests(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/pkg")).
		// Run tests
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v ./adapters)",
		})
}

func IntegrationTests(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/pkg")).
		// Dependencies
		With(ci.CockroachDependency(ci.CockroachDBService(client, "pkg"), "pkg")).
		With(ci.BinanceDependency(client)).
		// Run tests
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep ./adapters)",
		})
}
