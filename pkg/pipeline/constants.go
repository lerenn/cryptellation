package pipeline

var (
	ComponentsNames = []string{
		"backtests",
		"candlesticks",
		"exchanges",
		"indicators",
		"ticks",
	}
)

const (
	// CockroachDBImage is the image user for cockroach execution
	CockroachDBImage = "cockroachdb/cockroach"
	// GolangImage is the image used for golang execution.
	GolangImage = "golang:1.21.4"
	// LinterImage is the image used for linter.
	LinterImage = "golangci/golangci-lint:v1.55"
	// NATSImage is the image used for NATS.
	NATSImage = "nats:2.10"
	// RedisImage is the image used for Redis.
	RedisImage = "redis:6-alpine"

	// SecretsFilePath is the path to the file containing secrets
	SecretsFilePath = "./.credentials.env"

	// Source code path in go containers
	SourcePath = "/go/src/github.com/lerenn/cryptellation"
)
