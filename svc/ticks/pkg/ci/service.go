package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
	"github.com/lerenn/cryptellation/pkg/utils"
)

// BuildBinary builds the binaries for the service
func BuildBinary(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("golang:" + utils.GoVersion() + "-alpine3.19").
		// Add source code as work directory
		With(ci.SourceAsWorkdir(client, "/svc/"+ServiceName)).
		// Add command to generate code
		WithExec([]string{"go", "install", "./cmd/api", "./cmd/data"})
}

func Runner(client *dagger.Client) *dagger.Container {
	buildDir := BuildBinary(client).Directory("/go/bin")
	entrypointFile := client.Host().File("./svc/" + ServiceName + "/build/package/entrypoint.sh")

	return client.Container().
		// Add base image
		From("alpine:3.19").
		// Import binary
		WithDirectory("/usr/local/bin", buildDir).
		// Add entrypoint
		WithFile("/entrypoint.sh", entrypointFile).
		WithEntrypoint([]string{"sh", "/entrypoint.sh"}).
		// Add health port
		WithEnvVariable("HEALTH_PORT", "9000").
		WithExposedPort(9000).
		// Add default variables
		With(ci.DefaultSQLVariables()).
		With(ci.DefaultExchangesVariables(client)).
		With(ci.DefaultBrokerVariables())
}

func RunnerWithDependencies(client *dagger.Client, broker dagger.WithContainerFunc) *dagger.Container {
	return Runner(client).
		With(ci.CockroachDependency(ci.CockroachDBService(client, ServiceName), ServiceName)).
		With(ci.BinanceDependency(client)).
		With(broker)
}

func Service(client *dagger.Client, broker dagger.WithContainerFunc) *dagger.Service {
	return RunnerWithDependencies(client, broker).AsService()
}
