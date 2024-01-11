package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func Service(client *dagger.Client, broker dagger.WithContainerFunc) *dagger.Service {
	return client.Container().
		// Add base image
		From("golang:"+utils.GoVersion()).
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Add health port
		WithEnvVariable("HEALTH_PORT", "9000").
		WithExposedPort(9000).
		// Dependencies
		With(ci.CockroachDependency(ci.CockroachDB(client, ServiceName), ServiceName)).
		With(ci.BinanceDependency(client)).
		With(broker).
		// Run service with migrations
		WithExec([]string{"bash", "-c",
			"go run ./cmd/data migrations migrate && go run ./cmd/api serve",
		}).
		// As service
		AsService()
}
