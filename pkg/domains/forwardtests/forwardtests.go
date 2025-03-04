package forwardtests

import (
	"github.com/lerenn/cryptellation/v1/api"
	wfclient "github.com/lerenn/cryptellation/v1/clients/go/wfclient"
	"github.com/lerenn/cryptellation/v1/pkg/domains/forwardtests/activities/db"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ForwardTests is the interface for the forwardtests domain.
type ForwardTests interface {
	Register(w worker.Worker)

	CreateForwardTestWorkflow(
		ctx workflow.Context,
		params api.CreateForwardTestWorkflowParams,
	) (api.CreateForwardTestWorkflowResults, error)

	ListForwardTestsWorkflow(
		ctx workflow.Context,
		params api.ListForwardTestsWorkflowParams,
	) (api.ListForwardTestsWorkflowResults, error)

	CreateForwardTestOrderWorkflow(
		ctx workflow.Context,
		params api.CreateForwardTestOrderWorkflowParams,
	) (api.CreateForwardTestOrderWorkflowResults, error)

	ListForwardTestAccountsWorkflow(
		ctx workflow.Context,
		params api.ListForwardTestAccountsWorkflowParams,
	) (api.ListForwardTestAccountsWorkflowResults, error)

	GetForwardTestStatusWorkflow(
		ctx workflow.Context,
		params api.GetForwardTestStatusWorkflowParams,
	) (api.GetForwardTestStatusWorkflowResults, error)
}

var _ ForwardTests = &workflows{}

type workflows struct {
	db            db.DB
	cryptellation wfclient.Client
}

// New creates a new ForwardTests instance.
func New(db db.DB) ForwardTests {
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
	worker.RegisterWorkflowWithOptions(wf.CreateForwardTestWorkflow, workflow.RegisterOptions{
		Name: api.CreateForwardTestWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.ListForwardTestsWorkflow, workflow.RegisterOptions{
		Name: api.ListForwardTestsWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.CreateForwardTestOrderWorkflow, workflow.RegisterOptions{
		Name: api.CreateForwardTestOrderWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.ListForwardTestAccountsWorkflow, workflow.RegisterOptions{
		Name: api.ListForwardTestAccountsWorkflowName,
	})
	worker.RegisterWorkflowWithOptions(wf.GetForwardTestStatusWorkflow, workflow.RegisterOptions{
		Name: api.GetForwardTestStatusWorkflowName,
	})
}
