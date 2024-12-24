// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=db.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// CreateCandlesticksActivityName is the name of the CreateCandlesticks activity.
const CreateCandlesticksActivityName = "CreateCandlesticksActivity"

type (
	// CreateCandlesticksActivityParams is the parameters for the CreateCandlesticks activity.
	CreateCandlesticksActivityParams struct {
		List *candlestick.List
	}

	// CreateCandlesticksActivityResults is the result for the CreateCandlesticks activity.
	CreateCandlesticksActivityResults struct{}
)

// ReadCandlesticksActivityName is the name of the ReadCandlesticks activity.
const ReadCandlesticksActivityName = "ReadCandlesticksActivity"

type (
	// ReadCandlesticksActivityParams is the parameters for the ReadCandlesticks activity.
	ReadCandlesticksActivityParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    time.Time
		End      time.Time
		Limit    uint
	}

	// ReadCandlesticksActivityResults is the result for the ReadCandlesticks activity.
	ReadCandlesticksActivityResults struct {
		List *candlestick.List
	}
)

// UpdateCandlesticksActivityName is the name of the UpdateCandlesticks activity.
const UpdateCandlesticksActivityName = "UpdateCandlesticksActivity"

type (
	// UpdateCandlesticksActivityParams is the parameters for the UpdateCandlesticks activity.
	UpdateCandlesticksActivityParams struct {
		List *candlestick.List
	}

	// UpdateCandlesticksActivityResults is the result for the UpdateCandlesticks activity.
	UpdateCandlesticksActivityResults struct{}
)

// DeleteCandlesticksActivityName is the name of the DeleteCandlesticks activity.
const DeleteCandlesticksActivityName = "DeleteCandlesticksActivity"

type (
	// DeleteCandlesticksActivityParams is the parameters for the DeleteCandlesticks activity.
	DeleteCandlesticksActivityParams struct {
		List *candlestick.List
	}

	// DeleteCandlesticksActivityResults is the result for the DeleteCandlesticks activity.
	DeleteCandlesticksActivityResults struct{}
)

// DB is the interface that defines the candlesticks activities.
type DB interface {
	Register(w worker.Worker)

	CreateCandlesticksActivity(
		ctx context.Context,
		params CreateCandlesticksActivityParams,
	) (CreateCandlesticksActivityResults, error)

	ReadCandlesticksActivity(
		ctx context.Context,
		params ReadCandlesticksActivityParams,
	) (ReadCandlesticksActivityResults, error)

	UpdateCandlesticksActivity(
		ctx context.Context,
		params UpdateCandlesticksActivityParams,
	) (UpdateCandlesticksActivityResults, error)

	DeleteCandlesticksActivity(
		ctx context.Context,
		params DeleteCandlesticksActivityParams,
	) (DeleteCandlesticksActivityResults, error)
}

// DefaultActivityOptions returns the default database activities options.
func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{
				ErrNotFound.Error(),
			},
		},
		StartToCloseTimeout:    activities.DBStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.DBStartToCloseDefaultTimeout,
	}
}
