package docker

// RunnerInfo represents a Docker runner.
type RunnerInfo struct {
	OS              string
	Arch            string
	BuildBaseImage  string
	TargetBaseImage string
}

var (
	// GoRunnersInfo represents the different OS/Arch platform wanted for docker hub in Go service.
	GoRunnersInfo = map[string]RunnerInfo{
		"linux/386":      {OS: "linux", Arch: "386", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/amd64":    {OS: "linux", Arch: "amd64", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/arm/v6":   {OS: "linux", Arch: "arm/v6", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/arm/v7":   {OS: "linux", Arch: "arm/v7", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/arm64/v8": {OS: "linux", Arch: "arm64/v8", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/mips64le": {OS: "linux", Arch: "mips64le", BuildBaseImage: "golang", TargetBaseImage: "debian:12"},
		"linux/ppc64le":  {OS: "linux", Arch: "ppc64le", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		"linux/s390x":    {OS: "linux", Arch: "s390x", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
	}
)
