// A generated module for Cryptellation functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"maps"
	"runtime"
	"slices"

	"github.com/lerenn/cryptellation/v1/pkg/dagger/internal/dagger"
)

type Cryptellation struct{}

func (m *Cryptellation) AvailablePlatforms() []string {
	return slices.Collect(maps.Keys(GoRunnersInfo))
}

func (m *Cryptellation) Worker(
	sourceDir *dagger.Directory,
	// +optional
	targetPlatform string,
) *dagger.Container {
	// Get running OS, if that's an OS unsupported by Docker, replace by Linu
	os := runtime.GOOS
	if os == "darwin" {
		os = "linux"
	}

	// Set default runner info and override by argument
	runnerInfo := GoRunnersInfo["linux/amd64"]
	if targetPlatform != "" {
		info, ok := GoRunnersInfo[targetPlatform]
		if ok {
			runnerInfo = info
		}
	}

	return sourceDir.DockerBuild(dagger.DirectoryDockerBuildOpts{
		BuildArgs: []dagger.BuildArg{
			{Name: "BUILDPLATFORM", Value: os + "/" + runtime.GOARCH},
			{Name: "TARGETOS", Value: runnerInfo.OS},
			{Name: "TARGETARCH", Value: runnerInfo.Arch},
			{Name: "BUILDBASEIMAGE", Value: runnerInfo.BuildBaseImage},
			{Name: "TARGETBASEIMAGE", Value: runnerInfo.TargetBaseImage},
		},
		Platform:   dagger.Platform(runnerInfo.OS + "/" + runnerInfo.Arch),
		Dockerfile: "/build/package/Dockerfile",
	})
}

func (m *Cryptellation) WorkerWithDependencies(
	sourceDir *dagger.Directory,
	secretsFile *dagger.Secret,
	mongo *dagger.Service,
	nats *dagger.Service,
) *dagger.Container {
	c := m.Worker(sourceDir, runtime.GOOS+"/"+runtime.GOARCH)

	return c.WithExposedPort(9000, dagger.ContainerWithExposedPortOpts{
		Protocol:    dagger.NetworkProtocolTcp,
		Description: "Healthcheck",
	}).WithExec([]string{"worker", "serve"})
}
