package indicators

import (
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/clients/worker/go/wfclient"
	"github.com/lerenn/cryptellation/v1/pkg/domains/indicators/activities/db"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Indicators is the interface for the indicators domain.
type Indicators interface {
	Register(w worker.Worker)

	ListSMAWorkflow(
		ctx workflow.Context,
		params api.ListSMAWorkflowParams,
	) (api.ListSMAWorkflowResults, error)
}

// Check that the workflows implements the Indicators interface.
var _ Indicators = &workflows{}

type workflows struct {
	db            db.DB
	cryptellation wfclient.Client
}

// New creates a new Indicators instance.
func New(db db.DB) Indicators {
	if db == nil {
		panic("nil db")
	}

	return &workflows{
		cryptellation: wfclient.NewClient(),
		db:            db,
	}
}

// Register registers the workflows to the worker.
func (wf *workflows) Register(worker worker.Worker) {
	worker.RegisterWorkflowWithOptions(wf.ListSMAWorkflow, workflow.RegisterOptions{
		Name: api.ListSMAWorkflowName,
	})
}
