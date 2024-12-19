// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=exchanges.mock.gen.go -package exchanges

package exchanges

import (
	"context"
	"errors"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrInexistantExchange is the error when the exchange does not exist.
	ErrInexistantExchange = errors.New("inexistant exchange")
)

// ListenSymbolActivityName is the name of the activity to listen to a symbol.
const ListenSymbolActivityName = "ListenSymbolActivity"

type (
	// ListenSymbolParams is the parameters for the ListenSymbolActivity.
	ListenSymbolParams struct {
		ParentWorkflowID string
		Exchange         string
		Symbol           string
	}

	// ListenSymbolResults is the results for the ListenSymbolActivity.
	ListenSymbolResults struct{}
)

// Exchanges is the exchanges activities for ticks.
type Exchanges interface {
	Name() string
	Register(w worker.Worker)

	ListenSymbolActivity(ctx context.Context, params ListenSymbolParams) (ListenSymbolResults, error)
}

func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{},
		},
		StartToCloseTimeout:    activities.ExchangesStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.ExchangesScheduleToCloseDefaultTimeout,
	}
}
