package workflows

import (
	"fmt"

	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal"
	"go.temporal.io/sdk/workflow"
)

func (wf *workflows) ListenToTicksWorkflow(
	ctx workflow.Context,
	params internal.ListenToTicksWorkflowParams,
) (internal.ListenToTicksWorkflowResults, error) {
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	return internal.ListenToTicksWorkflowResults{}, nil
}
