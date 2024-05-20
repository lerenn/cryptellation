package ci

import (
	"dagger.io/dagger"
	"github.com/lerenn/cryptellation/pkg/ci"
)

const (
	// ServiceName is the name of the tested service
	ServiceName = "candlesticks"
)

func Runner(client *dagger.Client) *dagger.Container {
	return client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: "/svc/" + ServiceName + "/build/package/Dockerfile",
	})
}

func RunnerWithDependencies(client *dagger.Client, dependencies ...dagger.WithContainerFunc) *dagger.Container {
	r := Runner(client).
		With(ci.MongoDependency(ci.MongoService(client))).
		With(ci.BinanceDependency(client))

	for _, d := range dependencies {
		r = r.With(d)
	}

	return r
}

func Service(client *dagger.Client, broker dagger.WithContainerFunc) *dagger.Service {
	return RunnerWithDependencies(client, broker).AsService()
}
