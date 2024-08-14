package main

import "cryptellation/internal/dagger/internal/dagger"

func (mod *CryptellationInternal) Mongo() *dagger.Container {
	return dag.Container().
		// Add base image
		From("mongo:7-jammy").
		// Add exposed ports
		WithExposedPort(27017)
}

func (mod *CryptellationInternal) AttachMongo(c *dagger.Container, mongo *dagger.Service) *dagger.Container {
	return c.
		// Add service
		WithServiceBinding("mongo", mongo).
		// Add environment variables linked to service
		WithEnvVariable("MONGO_CONNECTION_STRING", "mongodb://mongo:27017").
		WithEnvVariable("MONGO_DATABASE", "cryptellation")
}
