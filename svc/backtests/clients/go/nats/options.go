package nats

import "github.com/lerenn/asyncapi-codegen/pkg/extensions"

type option func(b *nats)

func WithLogger(logger extensions.Logger) option {
	return func(b *nats) {
		b.logger = logger
	}
}

func WithName(name string) option {
	return func(b *nats) {
		b.name = name
	}
}
