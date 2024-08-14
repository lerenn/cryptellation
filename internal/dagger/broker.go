package main

import (
	"cryptellation/internal/dagger/internal/dagger"
)

// NATS creates a new NATS service.
func (mod *CryptellationInternal) NATS() *dagger.Container {
	return dag.Container().
		// Add base image
		From("nats:2.10").
		// Add exposed ports
		WithExposedPort(4222)
}

// AttachNATS attaches the NATS service to the container.
func (mod *CryptellationInternal) AttachNATS(c *dagger.Container, nats *dagger.Service) *dagger.Container {
	return c.
		// Add service
		WithServiceBinding("nats", nats).
		// Add environment variables linked to service
		WithEnvVariable("NATS_HOST", "nats").
		WithEnvVariable("NATS_PORT", "4222")
}
