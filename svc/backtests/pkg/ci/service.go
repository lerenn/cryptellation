package ci

import (
	"cryptellation/pkg/ci"

	"dagger.io/dagger"
)

const (
	// ServiceName is the name of the tested service
	ServiceName = "backtests"
)

func Runner(client *dagger.Client) *dagger.Container {
	return client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: "/svc/" + ServiceName + "/build/package/Dockerfile",
	})
}

func RunnerWithDependencies(client *dagger.Client, dependencies ...dagger.WithContainerFunc) *dagger.Container {
	r := Runner(client).
		With(ci.MongoDependency(ci.MongoService(client)))

	for _, d := range dependencies {
		r = r.With(d)
	}

	return r
}

func Service(client *dagger.Client, broker dagger.WithContainerFunc) *dagger.Service {
	return RunnerWithDependencies(client, broker).AsService()
}
