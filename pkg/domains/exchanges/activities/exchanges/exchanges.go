// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ListExchangesActivityName is the name of the ListExchanges activity.
const ListExchangesActivityName = "ListExchangesActivity"

type (
	// ListExchangesActivityParams is the parameters for the ListExchangesActivity activity.
	ListExchangesActivityParams struct{}

	// ListExchangesActivityResults is the result for the ListExchangesActivity activity.
	ListExchangesActivityResults struct {
		List []string
	}
)

// GetExchangeActivityName is the name of the GetExchange activity.
const GetExchangeActivityName = "GetExchangeActivity"

type (
	// GetExchangeActivityParams is the parameters for the GetExchange activity.
	GetExchangeActivityParams struct {
		Name string
	}

	// GetExchangeActivityResults is the result for the GetExchange activity.
	GetExchangeActivityResults struct {
		Exchange exchange.Exchange
	}
)

// Exchanges is the interface for the exchanges activities.
type Exchanges interface {
	Register(w worker.Worker)
	Name() string

	GetExchangeActivity(
		ctx context.Context,
		params GetExchangeActivityParams,
	) (GetExchangeActivityResults, error)

	ListExchangesActivity(
		ctx context.Context,
		params ListExchangesActivityParams,
	) (ListExchangesActivityResults, error)
}

func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		RetryPolicy: &temporal.RetryPolicy{
			NonRetryableErrorTypes: []string{
				ErrInexistantExchange.Error(),
			},
		},
		StartToCloseTimeout:    activities.ExchangesStartToCloseDefaultTimeout,
		ScheduleToCloseTimeout: activities.ExchangesScheduleToCloseDefaultTimeout,
	}
}
