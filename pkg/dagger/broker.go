package main

import (
	"cryptellation/pkg/dagger/internal/dagger"
)

// NATS creates a new NATS service.
func (mod *CryptellationPkg) NATS() *dagger.Container {
	return dag.Container().
		// Add base image
		From("nats:2.10").
		// Add exposed ports
		WithExposedPort(4222)
}

// AttachNATS attaches the NATS service to the container.
func (mod *CryptellationPkg) AttachNATS(c *dagger.Container, nats *dagger.Service) *dagger.Container {
	return c.
		// Add service
		WithServiceBinding("nats", nats).
		// Add environment variables linked to service
		WithEnvVariable("NATS_HOST", "nats").
		WithEnvVariable("NATS_PORT", "4222")
}
