package pipeline

import (
	"context"
	"fmt"
	"time"

	"dagger.io/dagger"
)

func ServiceCryptellation(client *dagger.Client, dependencies ...dagger.WithContainerFunc) map[string]*dagger.Service {
	services := make(map[string]*dagger.Service, len(ComponentsNames))
	for _, n := range ComponentsNames {
		container := containerWithSourceCodeAndGo(client).
			WithEnvVariable("SQLDB_DATABASE", n).
			WithEnvVariable("HEALTH_PORT", "9000").
			WithExposedPort(9000).
			With(DependsOnServices(dependencies...)).
			WithExec([]string{
				"go", "run", fmt.Sprintf("./cmd/cryptellation-%s", n), "serve",
			})

		services[n] = container.AsService()
	}
	return services
}

func DependsOnCryptellation(cryptellation map[string]*dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		for n, s := range cryptellation {
			r = r.WithServiceBinding("cryptellation-"+n, s)
		}
		return r
	}
}

func ServiceNATS(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Return container as a service
		AsService()
}

// DependsOnNATS returns a function that add a DependsOnNATS service to container
func DependsOnNATS(nats *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("nats", nats).
			// Add environment variables linked to service
			WithEnvVariable("NATS_HOST", "nats").
			WithEnvVariable("NATS_PORT", "4222")
	}
}

func ServiceRedis(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From(RedisImage).
		// Add exposed ports
		WithExposedPort(6379).
		// Return container as a service
		AsService()
}

// DependsOnRedis returns a function that add a DependsOnRedis service to container
func DependsOnRedis(redis *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("redis", redis).
			// Add environment variables linked to service
			WithEnvVariable("REDIS_URL", "redis:6379")
	}
}

func ServiceCockroachDB(ctx context.Context, client *dagger.Client) *dagger.Service {
	s := client.Container().
		// Add base image
		From(CockroachDBImage).
		// Add exposed ports
		WithExposedPort(26257).
		WithExposedPort(8080).
		// Add volumes
		WithDirectory("/docker-entrypoint-initdb.d", client.Host().Directory("./configs")).
		// Add command
		WithExec([]string{"start-single-node", "--insecure"}).
		// Return container as a service
		AsService()

	return s
}

func CockroachDBMigrations(client *dagger.Client, s *dagger.Service) []*dagger.Container {
	// Execute migrations
	migrations := make([]*dagger.Container, len(ComponentsNames))
	for i, n := range ComponentsNames {
		migrations[i] = containerWithSourceCodeAndGo(client).
			With(DependsOnCockroachDB(s)).
			WithEnvVariable("SQLDB_DATABASE", n).
			WithEnvVariable("TEST_TIME", time.Now().Format(time.RFC3339)). // Force replay
			WithExec([]string{
				"go", "run", fmt.Sprintf("./cmd/cryptellation-%s", n), "migrations", "migrate",
			})
	}
	return migrations
}

// DependsOnCockroachDB returns a function that add a DependsOnCockroachDB service to container
func DependsOnCockroachDB(s *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("cockroachdb", s).
			// Add environment variables linked to service
			WithEnvVariable("SQLDB_HOST", "cockroachdb").
			WithEnvVariable("SQLDB_PORT", "26257").
			WithEnvVariable("SQLDB_USER", "root").
			WithEnvVariable("SQLDB_PASSWORD", "")
	}
}

// DependsOnBinance returns a function that set variables to use binance as a service
func DependsOnBinance(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			With(secret(client, "BINANCE_API_KEY")).
			With(secret(client, "BINANCE_SECRET_KEY"))
	}
}

func DependsOnServices(services ...dagger.WithContainerFunc) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		for _, s := range services {
			r = r.With(s)
		}
		return r
	}
}
