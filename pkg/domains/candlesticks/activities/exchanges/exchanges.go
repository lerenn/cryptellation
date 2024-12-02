// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=mock.gen.go -package exchanges

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
	// GetCandlesticksActivityParams is the parameters for the GetCandlesticks activity.
	GetCandlesticksActivityParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    time.Time
		End      time.Time
		Limit    int
	}

	// GetCandlesticksActivityResults is the result for the GetCandlesticks activity.
	GetCandlesticksActivityResults struct {
		List *candlestick.List
	}
)

// Exchanges is the interface that defines the exchanges activities.
type Exchanges interface {
	Register(w worker.Worker)

	GetCandlesticksActivity(ctx context.Context, payload GetCandlesticksActivityParams) (GetCandlesticksActivityResults, error)
}
