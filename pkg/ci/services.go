package ci

import (
	"dagger.io/dagger"
)

func Nats(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From("nats:2.10").
		// Add exposed ports
		WithExposedPort(4222).
		// Return container as a service
		AsService()
}

// NatsDependency returns a function that add a NatsDependency service to container
func NatsDependency(nats *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("nats", nats).
			// Add environment variables linked to service
			WithEnvVariable("NATS_HOST", "nats").
			WithEnvVariable("NATS_PORT", "4222")
	}
}

func Redis(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From("redis:6-alpine").
		// Add exposed ports
		WithExposedPort(6379).
		// Return container as a service
		AsService()
}

// RedisDependency returns a function that add a DependsOnRedis service to container
func RedisDependency(redis *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("redis", redis).
			// Add environment variables linked to service
			WithEnvVariable("REDIS_URL", "redis:6379")
	}
}

func CockroachDB(client *dagger.Client, db string) *dagger.Service {
	s := client.Container().
		// Add base image
		From("cockroachdb/cockroach").
		// Add exposed ports
		WithExposedPort(26257).
		WithExposedPort(8080).
		// With database
		WithEnvVariable("COCKROACH_DATABASE", db).
		// Add command
		WithExec([]string{"start-single-node", "--insecure"}).
		// Return container as a service
		AsService()

	return s
}

// CockroachDependency returns a function that add a CockroachDependency service to container
func CockroachDependency(s *dagger.Service, db string) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("cockroachdb", s).
			// Add environment variables linked to service
			WithEnvVariable("SQLDB_HOST", "cockroachdb").
			WithEnvVariable("SQLDB_PORT", "26257").
			WithEnvVariable("SQLDB_USER", "root").
			WithEnvVariable("SQLDB_PASSWORD", "").
			WithEnvVariable("SQLDB_DATABASE", db)
	}
}

// BinanceDependency returns a function that set variables to use binance as a service
func BinanceDependency(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			With(Secret(client, "BINANCE_API_KEY")).
			With(Secret(client, "BINANCE_SECRET_KEY"))
	}
}
