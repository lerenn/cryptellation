package temporal

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
)

// CreateForwardtestWorkflowName is the name of the CreateForwardtestWorkflow.
const CreateForwardtestWorkflowName = "CreateForwardtestWorkflow"

type (
	// CreateForwardtestWorkflowParams is the input for the CreateForwardtestWorkflow.
	CreateForwardtestWorkflowParams struct {
		Accounts map[string]account.Account
	}

	// CreateForwardtestWorkflowResults is the output for the CreateForwardtestWorkflow.
	CreateForwardtestWorkflowResults struct {
		ID uuid.UUID
	}
)

// ListForwardtestsWorkflowName is the name of the ListForwardtestsWorkflow.
const ListForwardtestsWorkflowName = "ListForwardtestsWorkflow"

type (
	// ListForwardtestsWorkflowParams is the input for the ListForwardtestsWorkflow.
	ListForwardtestsWorkflowParams struct{}

	// ListForwardtestsWorkflowResults is the output for the ListForwardtestsWorkflow.
	ListForwardtestsWorkflowResults struct {
		Forwardtests []forwardtest.Forwardtest
	}
)

// CreateForwardtestOrderWorkflowName is the name of the CreateForwardtestOrderWorkflow.
const CreateForwardtestOrderWorkflowName = "CreateForwardtestOrderWorkflow"

type (
	// CreateForwardtestOrderWorkflowParams is the input for the CreateForwardtestOrderWorkflow.
	CreateForwardtestOrderWorkflowParams struct {
		ForwardtestID uuid.UUID
		Order         order.Order
	}

	// CreateForwardtestOrderWorkflowResults is the output for the CreateForwardtestOrderWorkflow.
	CreateForwardtestOrderWorkflowResults struct{}
)

// ListForwardtestAccountsWorkflowName is the name of the ListForwardtestAccountsWorkflow.
const ListForwardtestAccountsWorkflowName = "ListForwardtestAccountsWorkflow"

type (
	// ListForwardtestAccountsWorkflowParams is the input for the ListForwardtestAccountsWorkflow.
	ListForwardtestAccountsWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// ListForwardtestAccountsWorkflowResults is the output for the ListForwardtestAccountsWorkflow.
	ListForwardtestAccountsWorkflowResults struct {
		Accounts map[string]account.Account
	}
)

// GetForwardtestStatusWorkflowName is the name of the GetForwardtestStatusWorkflow.
const GetForwardtestStatusWorkflowName = "GetForwardtestStatusWorkflow"

type (
	// GetForwardtestStatusWorkflowParams is the input for the GetForwardtestStatusWorkflow.
	GetForwardtestStatusWorkflowParams struct {
		ForwardtestID uuid.UUID
	}

	// GetForwardtestStatusWorkflowResults is the output for the GetForwardtestStatusWorkflow.
	GetForwardtestStatusWorkflowResults struct {
		Status forwardtest.Status
	}
)
