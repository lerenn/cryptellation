// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=interface.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"go.temporal.io/sdk/worker"
)

const (
	// GetCandlesticksActivityName is the name of the GetCandlesticks activity.
	GetCandlesticksActivityName = "GetCandlesticksActivity"
)

type (
	// GetCandlesticksParams is the parameters for the GetCandlesticks activity.
	GetCandlesticksParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    time.Time
		End      time.Time
		Limit    int
	}

	// GetCandlesticksResult is the result for the GetCandlesticks activity.
	GetCandlesticksResult struct {
		List *candlestick.List
	}
)

// Interface is the interface that defines the GetCandlesticks activity.
type Interface interface {
	Register(w worker.Worker)

	GetCandlesticks(ctx context.Context, payload GetCandlesticksParams) (GetCandlesticksResult, error)
}
