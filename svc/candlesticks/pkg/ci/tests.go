package ci

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
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Test
		WithExec([]string{"sh", "-c",
			"go test $(go list ./... | grep -v -e ./internal/adapters -e ./test)",
		})
}

func IntegrationTests(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Dependencies
		With(ci.MongoDependency(ci.MongoService(client))).
		With(ci.BinanceDependency(client)).
		// Run tests
		WithExec([]string{"sh", "-c",
			"go test ./internal/adapters/...",
		})
}

func EndToEndTests(client *dagger.Client) *dagger.Container {
	broker := ci.NatsService(client)
	service := Service(client, ci.NatsDependency(broker))

	return client.Container().
		// Add base image
		From("golang:"+utils.GoVersion()+"-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Dependencies
		WithServiceBinding("cryptellation-candlesticks", service).
		With(ci.NatsDependency(broker)).
		// Run tests
		WithExec([]string{
			"go", "test", "./test/...",
		})
}
