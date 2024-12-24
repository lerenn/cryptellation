package candlesticks

import (
	"errors"
	"fmt"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/db"
	"github.com/lerenn/cryptellation/v1/pkg/domains/candlesticks/activities/exchanges"
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

	ListCandlesticksWorkflow(
		ctx workflow.Context,
		payload api.ListCandlesticksWorkflowParams,
	) (api.ListCandlesticksWorkflowResults, error)
}

// Check that the workflows implements the Candlesticks interface.
var _ Candlesticks = &workflows{}

type workflows struct {
	db        db.DB
	exchanges exchanges.Exchanges
}

// New creates a new candlesticks workflows.
func New(db db.DB, exchanges exchanges.Exchanges) Candlesticks {
	if db == nil {
		panic("nil db")
	}

	if exchanges == nil {
		panic("nil exchanges")
	}

	return &workflows{
		db:        db,
		exchanges: exchanges,
	}
}

// Register registers the candlesticks workflows to the worker.
func (wf *workflows) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(wf.ListCandlesticksWorkflow, workflow.RegisterOptions{
		Name: api.ListCandlesticksWorkflowName,
	})
}
