// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=db.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// CreateForwardtestActivityName is the name of the CreateForwardtestActivity.
const CreateForwardtestActivityName = "CreateForwardtestActivity"

type (
	// CreateForwardtestActivityParams is the parameters for the CreateForwardtestActivity.
	CreateForwardtestActivityParams struct {
		Forwardtest forwardtest.Forwardtest
	}

	// CreateForwardtestActivityResult is the result for the CreateForwardtestActivity.
	CreateForwardtestActivityResult struct{}
)

// ReadForwardtestActivityName is the name of the ReadForwardtestActivity.
const ReadForwardtestActivityName = "ReadForwardtestActivity"

type (
	// ReadForwardtestActivityParams is the parameters for the ReadForwardtestActivity.
	ReadForwardtestActivityParams struct {
		ID uuid.UUID
	}

	// ReadForwardtestActivityResult is the result for the ReadForwardtestActivity.
	ReadForwardtestActivityResult struct {
		Forwardtest forwardtest.Forwardtest
	}
)

// ListForwardtestsActivityName is the name of the ListForwardtestsActivity.
const ListForwardtestsActivityName = "ListForwardtestsActivity"

type (
	// ListForwardtestsActivityParams is the parameters for the ListForwardtestsActivity.
	ListForwardtestsActivityParams struct{}

	// ListForwardtestsActivityResult is the result for the ListForwardtestsActivity.
	ListForwardtestsActivityResult struct {
		Forwardtests []forwardtest.Forwardtest
	}
)

// UpdateForwardtestActivityName is the name of the UpdateForwardtestActivity.
const UpdateForwardtestActivityName = "UpdateForwardtestActivity"

type (
	// UpdateForwardtestActivityParams is the parameters for the UpdateForwardtestActivity.
	UpdateForwardtestActivityParams struct {
		Forwardtest forwardtest.Forwardtest
	}

	// UpdateForwardtestActivityResult is the result for the UpdateForwardtestActivity.
	UpdateForwardtestActivityResult struct{}
)

// DeleteForwardtestActivityName is the name of the DeleteForwardtestActivity.
const DeleteForwardtestActivityName = "DeleteForwardtestActivity"

type (
	// DeleteForwardtestActivityParams is the parameters for the DeleteForwardtestActivity.
	DeleteForwardtestActivityParams struct {
		ID uuid.UUID
	}

	// DeleteForwardtestActivityResult is the result for the DeleteForwardtestActivity.
	DeleteForwardtestActivityResult struct{}
)

// DB is the interface for the database activities.
type DB interface {
	Register(w worker.Worker)

	CreateForwardtestActivity(
		ctx context.Context,
		params CreateForwardtestActivityParams,
	) (CreateForwardtestActivityResult, error)
	ReadForwardtestActivity(
		ctx context.Context,
		params ReadForwardtestActivityParams,
	) (ReadForwardtestActivityResult, error)
	ListForwardtestsActivity(
		ctx context.Context,
		params ListForwardtestsActivityParams,
	) (ListForwardtestsActivityResult, error)
	UpdateForwardtestActivity(
		ctx context.Context,
		params UpdateForwardtestActivityParams,
	) (UpdateForwardtestActivityResult, error)
	DeleteForwardtestActivity(
		ctx context.Context,
		params DeleteForwardtestActivityParams,
	) (DeleteForwardtestActivityResult, error)
}

// DefaultActivityOptions returns the default database activities options.
func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{
				ErrRecordNotFound.Error(),
				ErrNotImplemented.Error(),
				ErrNilID.Error(),
			},
		},
		StartToCloseTimeout:    activities.DBStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.DBStartToCloseDefaultTimeout,
	}
}
