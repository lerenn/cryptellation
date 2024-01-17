package ci

import (
	"context"

	"dagger.io/dagger"
)

// PublishDockerImage builds the docker image for the service
func PublishDockerImage(client *dagger.Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, err := Runner(client).
			Publish(ctx, "lerenn/cryptellation-"+ServiceName+":latest")
		return err
	}
}
