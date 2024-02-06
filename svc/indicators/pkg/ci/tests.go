package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/utils"
	candlesticksCi "github.com/lerenn/cryptellation/svc/candlesticks/pkg/ci"
	"github.com/lerenn/cryptellation/tools/pkg/ci"
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
		With(ci.CockroachDependency(ci.CockroachDBService(client, ServiceName), ServiceName)).
		// Run tests
		WithExec([]string{"sh", "-c",
			"go run ./cmd/data migrations migrate && go test ./internal/adapters/...",
		})
}

func EndToEndTests(client *dagger.Client) *dagger.Container {
	broker := ci.NatsService(client)
	candlesticks := candlesticksCi.Service(client, ci.NatsDependency(broker))
	service := Service(client, ci.NatsDependency(broker), candlesticks)

	return client.Container().
		// Add base image
		From("golang:"+utils.GoVersion()+"-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Dependencies
		WithServiceBinding("cryptellation", service).
		With(ci.NatsDependency(broker)).
		// Run tests
		WithExec([]string{
			"go", "test", "./test/...",
		})
}
