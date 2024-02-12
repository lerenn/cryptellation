package ci

import (
	"dagger.io/dagger"
)

func Nats(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("nats:2.10").
		// Add exposed ports
		WithExposedPort(4222)
}

func NatsService(client *dagger.Client) *dagger.Service {
	return Nats(client).AsService()
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
