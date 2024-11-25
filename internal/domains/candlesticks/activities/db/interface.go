// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=interface.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"go.temporal.io/sdk/worker"
)

// CreateCandlesticksActivityName is the name of the CreateCandlesticks activity.
const CreateCandlesticksActivityName = "CreateCandlesticksActivity"

type (
	// CreateCandlesticksParams is the parameters for the CreateCandlesticks activity.
	CreateCandlesticksParams struct {
		List *candlestick.List
	}

	// CreateCandlesticksResult is the result for the CreateCandlesticks activity.
	CreateCandlesticksResult struct{}
)

// ReadCandlesticksActivityName is the name of the ReadCandlesticks activity.
const ReadCandlesticksActivityName = "ReadCandlesticksActivity"

type (
	// ReadCandlesticksParams is the parameters for the ReadCandlesticks activity.
	ReadCandlesticksParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    time.Time
		End      time.Time
		Limit    uint
	}

	// ReadCandlesticksResult is the result for the ReadCandlesticks activity.
	ReadCandlesticksResult struct {
		List *candlestick.List
	}
)

// UpdateCandlesticksActivityName is the name of the UpdateCandlesticks activity.
const UpdateCandlesticksActivityName = "UpdateCandlesticksActivity"

type (
	// UpdateCandlesticksParams is the parameters for the UpdateCandlesticks activity.
	UpdateCandlesticksParams struct {
		List *candlestick.List
	}

	// UpdateCandlesticksResult is the result for the UpdateCandlesticks activity.
	UpdateCandlesticksResult struct{}
)

// DeleteCandlesticksActivityName is the name of the DeleteCandlesticks activity.
const DeleteCandlesticksActivityName = "DeleteCandlesticksActivity"

type (
	// DeleteCandlesticksParams is the parameters for the DeleteCandlesticks activity.
	DeleteCandlesticksParams struct {
		List *candlestick.List
	}

	// DeleteCandlesticksResult is the result for the DeleteCandlesticks activity.
	DeleteCandlesticksResult struct{}
)

// Interface is the interface that defines the candlesticks activities.
type Interface interface {
	Register(w worker.Worker)

	CreateCandlesticks(ctx context.Context, params CreateCandlesticksParams) (CreateCandlesticksResult, error)
	ReadCandlesticks(ctx context.Context, params ReadCandlesticksParams) (ReadCandlesticksResult, error)
	UpdateCandlesticks(ctx context.Context, params UpdateCandlesticksParams) (UpdateCandlesticksResult, error)
	DeleteCandlesticks(ctx context.Context, params DeleteCandlesticksParams) (DeleteCandlesticksResult, error)
}
