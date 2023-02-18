package candlesticks

import "github.com/digital-feather/cryptellation/pkg/version"

var (
	Version = version.Version{
		SemVer:     "1.0.0",
		CommitHash: version.DefaultHash,
	}
)
