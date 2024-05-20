package ci

import "dagger.io/dagger"

func Mongo(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From("mongo:7-jammy").
		// Add exposed ports
		WithExposedPort(27017)
}

func MongoService(client *dagger.Client) *dagger.Service {
	return Mongo(client).AsService()
}

// MongoDependency returns a function that add a DependsOnMongo service to container
func MongoDependency(mongo *dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add service
			WithServiceBinding("mongo", mongo).
			// Add environment variables linked to service
			WithEnvVariable("MONGO_CONNECTION_STRING", "mongodb://mongo:27017")
	}
}
