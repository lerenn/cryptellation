// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=db.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/backtest"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// CreateBacktestActivityName is the name of the activity to create a backtest.
const CreateBacktestActivityName = "CreateBacktestActivity"

type (
	// CreateBacktestActivityParams is the parameters of the CreateBacktestActivity activity.
	CreateBacktestActivityParams struct {
		Backtest backtest.Backtest
	}

	// CreateBacktestActivityResults is the results of the CreateBacktestActivity activity.
	CreateBacktestActivityResults struct{}
)

// ReadBacktestActivityName is the name of the activity to read a backtest.
const ReadBacktestActivityName = "ReadBacktestActivity"

type (
	// ReadBacktestActivityParams is the parameters of the ReadBacktestActivity activity.
	ReadBacktestActivityParams struct {
		ID uuid.UUID
	}

	// ReadBacktestActivityResults is the results of the ReadBacktestActivity activity.
	ReadBacktestActivityResults struct {
		Backtest backtest.Backtest
	}
)

// ListBacktestsActivityName is the name of the activity to list backtests.
const ListBacktestsActivityName = "ListBacktestsActivity"

type (
	// ListBacktestsActivityParams is the parameters of the ListBacktestsActivity activity.
	ListBacktestsActivityParams struct{}

	// ListBacktestsActivityResults is the results of the ListBacktestsActivity activity.
	ListBacktestsActivityResults struct {
		Backtests []backtest.Backtest
	}
)

// UpdateBacktestActivityName is the name of the activity to update a backtest.
const UpdateBacktestActivityName = "UpdateBacktestActivity"

type (
	// UpdateBacktestActivityParams is the parameters of the UpdateBacktestActivity activity.
	UpdateBacktestActivityParams struct {
		Backtest backtest.Backtest
	}

	// UpdateBacktestActivityResults is the results of the UpdateBacktestActivity activity.
	UpdateBacktestActivityResults struct{}
)

// DeleteBacktestActivityName is the name of the activity to delete a backtest.
const DeleteBacktestActivityName = "DeleteBacktestActivity"

type (
	// DeleteBacktestActivityParams is the parameters of the DeleteBacktestActivity activity.
	DeleteBacktestActivityParams struct {
		ID uuid.UUID
	}

	// DeleteBacktestActivityResults is the results of the DeleteBacktestActivity activity.
	DeleteBacktestActivityResults struct{}
)

// DB is the interface for the backtest activity database.
type DB interface {
	Register(w worker.Worker)

	CreateBacktestActivity(
		ctx context.Context,
		params CreateBacktestActivityParams,
	) (CreateBacktestActivityResults, error)
	ReadBacktestActivity(
		ctx context.Context,
		params ReadBacktestActivityParams,
	) (ReadBacktestActivityResults, error)
	ListBacktestsActivity(
		ctx context.Context,
		params ListBacktestsActivityParams,
	) (ListBacktestsActivityResults, error)
	UpdateBacktestActivity(
		ctx context.Context,
		params UpdateBacktestActivityParams,
	) (UpdateBacktestActivityResults, error)
	DeleteBacktestActivity(
		ctx context.Context,
		params DeleteBacktestActivityParams,
	) (DeleteBacktestActivityResults, error)
}

// DefaultActivityOptions returns the default database activities options.
func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{
				ErrNilID.Error(),
				ErrNotFound.Error(),
				ErrNotImplemented.Error(),
			},
		},
		StartToCloseTimeout:    activities.DBStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.DBStartToCloseDefaultTimeout,
	}
}
