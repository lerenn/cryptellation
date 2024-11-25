package workflows

import (
	"errors"
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrCandlesticksWorkflow is returned when an error occurs in the candlesticks workflow.
	ErrCandlesticksWorkflow = errors.New("error during candlesticks workflow")
	// ErrNoExchange is returned when no exchange is found.
	ErrNoExchange = fmt.Errorf("%w: no exchange", ErrCandlesticksWorkflow)
	// ErrNoPair is returned when no pair is found.
	ErrNoPair = fmt.Errorf("%w: no pair", ErrCandlesticksWorkflow)
	// ErrNoPeriod is returned when no period is found.
	ErrNoPeriod = fmt.Errorf("%w: no period", ErrCandlesticksWorkflow)
)

// Candlesticks is the interface that describe the candlesticks workflows.
type Candlesticks interface {
	Register(w worker.Worker)

	ListCandlesticks(ctx workflow.Context, payload api.ListCandlesticksParams) (api.ListCandlesticksResults, error)
}
