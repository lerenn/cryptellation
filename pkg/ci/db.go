package ci

import "dagger.io/dagger"

func Redis(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("redis:6-alpine").
		// Add exposed ports
		WithExposedPort(6379)
}

func RedisService(client *dagger.Client) *dagger.Service {
	return Redis(client).AsService()
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

func CockroachDB(client *dagger.Client, db string) *dagger.Container {
	return client.Container().
		// Add base image
		From("cockroachdb/cockroach").
		// Add exposed ports
		WithExposedPort(26257).
		WithExposedPort(8080).
		// With database
		WithEnvVariable("COCKROACH_DATABASE", db).
		// Add command
		WithExec([]string{"start-single-node", "--insecure"})
}

func CockroachDBService(client *dagger.Client, db string) *dagger.Service {
	return CockroachDB(client, db).AsService()
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
