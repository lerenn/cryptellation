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

// CreateForwardTestActivityName is the name of the CreateForwardTestActivity.
const CreateForwardTestActivityName = "CreateForwardTestActivity"

type (
	// CreateForwardTestActivityParams is the parameters for the CreateForwardTestActivity.
	CreateForwardTestActivityParams struct {
		ForwardTest forwardtest.ForwardTest
	}

	// CreateForwardTestActivityResult is the result for the CreateForwardTestActivity.
	CreateForwardTestActivityResult struct{}
)

// ReadForwardTestActivityName is the name of the ReadForwardTestActivity.
const ReadForwardTestActivityName = "ReadForwardTestActivity"

type (
	// ReadForwardTestActivityParams is the parameters for the ReadForwardTestActivity.
	ReadForwardTestActivityParams struct {
		ID uuid.UUID
	}

	// ReadForwardTestActivityResult is the result for the ReadForwardTestActivity.
	ReadForwardTestActivityResult struct {
		ForwardTest forwardtest.ForwardTest
	}
)

// ListForwardTestsActivityName is the name of the ListForwardTestsActivity.
const ListForwardTestsActivityName = "ListForwardTestsActivity"

type (
	// ListForwardTestsActivityParams is the parameters for the ListForwardTestsActivity.
	ListForwardTestsActivityParams struct{}

	// ListForwardTestsActivityResult is the result for the ListForwardTestsActivity.
	ListForwardTestsActivityResult struct {
		ForwardTests []forwardtest.ForwardTest
	}
)

// UpdateForwardTestActivityName is the name of the UpdateForwardTestActivity.
const UpdateForwardTestActivityName = "UpdateForwardTestActivity"

type (
	// UpdateForwardTestActivityParams is the parameters for the UpdateForwardTestActivity.
	UpdateForwardTestActivityParams struct {
		ForwardTest forwardtest.ForwardTest
	}

	// UpdateForwardTestActivityResult is the result for the UpdateForwardTestActivity.
	UpdateForwardTestActivityResult struct{}
)

// DeleteForwardTestActivityName is the name of the DeleteForwardTestActivity.
const DeleteForwardTestActivityName = "DeleteForwardTestActivity"

type (
	// DeleteForwardTestActivityParams is the parameters for the DeleteForwardTestActivity.
	DeleteForwardTestActivityParams struct {
		ID uuid.UUID
	}

	// DeleteForwardTestActivityResult is the result for the DeleteForwardTestActivity.
	DeleteForwardTestActivityResult struct{}
)

// DB is the interface for the database activities.
type DB interface {
	Register(w worker.Worker)

	CreateForwardTestActivity(
		ctx context.Context,
		params CreateForwardTestActivityParams,
	) (CreateForwardTestActivityResult, error)
	ReadForwardTestActivity(
		ctx context.Context,
		params ReadForwardTestActivityParams,
	) (ReadForwardTestActivityResult, error)
	ListForwardTestsActivity(
		ctx context.Context,
		params ListForwardTestsActivityParams,
	) (ListForwardTestsActivityResult, error)
	UpdateForwardTestActivity(
		ctx context.Context,
		params UpdateForwardTestActivityParams,
	) (UpdateForwardTestActivityResult, error)
	DeleteForwardTestActivity(
		ctx context.Context,
		params DeleteForwardTestActivityParams,
	) (DeleteForwardTestActivityResult, error)
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
