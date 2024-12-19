// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=db.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// CreateExchangesActivityName is the name of the CreateExchanges activity.
const CreateExchangesActivityName = "CreateExchangesActivity"

type (
	// CreateExchangesActivityParams is the parameters for the CreateExchanges activity.
	CreateExchangesActivityParams struct {
		Exchanges []exchange.Exchange
	}

	// CreateExchangesActivityResults is the result for the CreateExchanges activity.
	CreateExchangesActivityResults struct{}
)

// ReadExchangesActivityName is the name of the ReadExchanges activity.
const ReadExchangesActivityName = "ReadExchangesActivity"

type (
	// ReadExchangesActivityParams is the parameters for the ReadExchanges activity.
	ReadExchangesActivityParams struct {
		Names []string
	}

	// ReadExchangesActivityResults is the result for the ReadExchanges activity.
	ReadExchangesActivityResults struct {
		Exchanges []exchange.Exchange
	}
)

// UpdateExchangesActivityName is the name of the UpdateExchanges activity.
const UpdateExchangesActivityName = "UpdateExchangesActivity"

type (
	// UpdateExchangesActivityParams is the parameters for the UpdateExchanges activity.
	UpdateExchangesActivityParams struct {
		Exchanges []exchange.Exchange
	}

	// UpdateExchangesActivityResults is the result for the UpdateExchanges activity.
	UpdateExchangesActivityResults struct{}
)

// DeleteExchangesActivityName is the name of the DeleteExchanges activity.
const DeleteExchangesActivityName = "DeleteExchangesActivity"

type (
	// DeleteExchangesActivityParams is the parameters for the DeleteExchanges activity.
	DeleteExchangesActivityParams struct {
		Names []string
	}

	// DeleteExchangesActivityResults is the result for the DeleteExchanges activity.
	DeleteExchangesActivityResults struct{}
)

// DB is the interface that the database activities must implement.
type DB interface {
	Register(w worker.Worker)

	CreateExchangesActivity(
		ctx context.Context,
		params CreateExchangesActivityParams,
	) (CreateExchangesActivityResults, error)

	ReadExchangesActivity(
		ctx context.Context,
		params ReadExchangesActivityParams,
	) (ReadExchangesActivityResults, error)

	UpdateExchangesActivity(
		ctx context.Context,
		params UpdateExchangesActivityParams,
	) (UpdateExchangesActivityResults, error)

	DeleteExchangesActivity(
		ctx context.Context,
		params DeleteExchangesActivityParams,
	) (DeleteExchangesActivityResults, error)
}

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
