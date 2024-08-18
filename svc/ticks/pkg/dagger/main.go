// A generated module for CryptellationTicks functions
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
	"cryptellation/internal/docker"
	"cryptellation/svc/ticks/pkg/dagger/internal/dagger"
	"runtime"
)

type CryptellationTicks struct{}

func (m *CryptellationTicks) Runner(
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
	runnerInfo := docker.GoRunnersInfo["linux/amd64"]
	if targetPlatform != "" {
		info, ok := docker.GoRunnersInfo[targetPlatform]
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
		Dockerfile: "/svc/ticks/build/package/Dockerfile",
	})
}

func (m *CryptellationTicks) RunnerWithDependencies(
	sourceDir *dagger.Directory,
	binanceSecretsFile *dagger.Secret,
	mongoService *dagger.Service,
	natsService *dagger.Service,
) *dagger.Container {
	c := m.Runner(sourceDir, runtime.GOOS+"/"+runtime.GOARCH)

	c = dag.CryptellationInternal().AttachMongo(c, mongoService)
	c = dag.CryptellationInternal().AttachBinance(c, binanceSecretsFile)
	c = dag.CryptellationInternal().AttachNats(c, natsService)

	return c.WithExposedPort(9000, dagger.ContainerWithExposedPortOpts{
		Protocol:    dagger.Tcp,
		Description: "Healthcheck",
	}).WithExec([]string{"api", "serve"})
}
