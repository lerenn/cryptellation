package pipeline

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

func UnitTests(client *dagger.Client) []*dagger.Container {
	return []*dagger.Container{
		// Commands
		containerWithSourceCodeAndGo(client).
			WithExec([]string{"go", "test", "./cmd/..."}),

		// Packages
		containerWithSourceCodeAndGo(client).
			WithExec([]string{"go", "test", "./pkg/..."}),

		// Components
		containerWithSourceCodeAndGo(client).
			WithExec([]string{"go", "test", "./internal/components/..."}),

		// Controllers
		containerWithSourceCodeAndGo(client).
			WithExec([]string{"go", "test", "./internal/controllers/..."}),
	}
}

func IntegrationTests(ctx context.Context, client *dagger.Client) (preflights, tests []*dagger.Container) {
	sqlService := ServiceCockroachDB(ctx, client)

	preflights = make([]*dagger.Container, 0)
	preflights = append(preflights, CockroachDBMigrations(client, sqlService)...)

	tests = []*dagger.Container{
		// NATS
		containerWithSourceCodeAndGo(client).
			With(DependsOnNATS(ServiceNATS(client))).
			WithExec([]string{"go", "test", "./internal/adapters/events/nats/..."}),

		// Redis
		containerWithSourceCodeAndGo(client).
			With(DependsOnRedis(ServiceRedis(client))).
			WithExec([]string{"go", "test", "./internal/adapters/db/redis/..."}),

		// SQL
		containerWithSourceCodeAndGo(client).
			With(DependsOnCockroachDB(sqlService)).
			WithExec([]string{"go", "test", "./internal/adapters/db/sql/..."}),

		// Binance
		containerWithSourceCodeAndGo(client).
			With(DependsOnBinance(client)).
			WithExec([]string{"go", "test", "./internal/adapters/exchanges/binance/..."}),
	}

	return preflights, tests
}

func EndToEndTests(ctx context.Context, client *dagger.Client) (preflights, tests []*dagger.Container) {
	// Generate preflights neede by tests
	sqlService := ServiceCockroachDB(ctx, client)
	preflights = make([]*dagger.Container, 0)
	preflights = append(preflights, CockroachDBMigrations(client, sqlService)...)

	// Generate services needed by tests
	services := []dagger.WithContainerFunc{
		DependsOnNATS(ServiceNATS(client)),
		DependsOnCockroachDB(sqlService),
		DependsOnBinance(client),
	}
	cryptellation := ServiceCryptellation(client, services...)
	services = append(services, DependsOnCryptellation(cryptellation))

	// Generate tests
	tests = make([]*dagger.Container, len(ComponentsNames))
	for i, n := range ComponentsNames {
		c := containerWithSourceCodeAndGo(client).
			With(DependsOnServices(services...)).
			WithExec([]string{"go", "test", fmt.Sprintf("./test/%s_test.go", n)})
		tests[i] = c
	}

	return preflights, tests
}
