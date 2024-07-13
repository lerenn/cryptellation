package nats

import "github.com/lerenn/asyncapi-codegen/pkg/extensions"

type option func(i *nats)

func WithLogger(logger extensions.Logger) option {
	return func(c *nats) {
		c.logger = logger
	}
}

func WithName(name string) option {
	return func(c *nats) {
		c.name = name
	}
}
