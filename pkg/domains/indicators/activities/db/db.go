// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=db.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ReadSMAActivityName is the name of the GetSMA activity.
const ReadSMAActivityName = "ReadSMAActivity"

type (
	// ReadSMAActivityParams is the parameters for the GetSMA activity.
	ReadSMAActivityParams struct {
		Exchange     string
		Pair         string
		Period       period.Symbol
		PeriodNumber int
		PriceType    candlestick.PriceType
		Start        time.Time
		End          time.Time
	}

	// ReadSMAActivityResults is the result for the GetSMA activity.
	ReadSMAActivityResults struct {
		Data *timeserie.TimeSerie[float64]
	}
)

// UpsertSMAActivityName is the name of the UpsertSMA activity.
const UpsertSMAActivityName = "UpsertSMAActivity"

type (
	// UpsertSMAActivityParams is the parameters for the UpsertSMA activity.
	UpsertSMAActivityParams struct {
		Exchange     string
		Pair         string
		Period       period.Symbol
		PeriodNumber int
		PriceType    candlestick.PriceType
		TimeSerie    *timeserie.TimeSerie[float64]
	}

	// UpsertSMAActivityResults is the result for the UpsertSMA activity.
	UpsertSMAActivityResults struct{}
)

// DB is the interface for the database activities.
type DB interface {
	Register(w worker.Worker)

	ReadSMAActivity(
		ctx context.Context,
		params ReadSMAActivityParams,
	) (ReadSMAActivityResults, error)

	UpsertSMAActivity(
		ctx context.Context,
		params UpsertSMAActivityParams,
	) (UpsertSMAActivityResults, error)
}

// DefaultActivityOptions returns the default database activities options.
func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{
				ErrNilID.Error(),
				ErrNoDocument.Error(),
			},
		},
		StartToCloseTimeout:    activities.DBStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.DBStartToCloseDefaultTimeout,
	}
}
