package api

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
)

// CreateForwardTestWorkflowName is the name of the CreateForwardTestWorkflow.
const CreateForwardTestWorkflowName = "CreateForwardTestWorkflow"

type (
	// CreateForwardTestWorkflowParams is the input for the CreateForwardTestWorkflow.
	CreateForwardTestWorkflowParams struct {
		Accounts map[string]account.Account
	}

	// CreateForwardTestWorkflowResults is the output for the CreateForwardTestWorkflow.
	CreateForwardTestWorkflowResults struct {
		ID uuid.UUID
	}
)

// ListForwardTestsWorkflowName is the name of the ListForwardTestsWorkflow.
const ListForwardTestsWorkflowName = "ListForwardTestsWorkflow"

type (
	// ListForwardTestsWorkflowParams is the input for the ListForwardTestsWorkflow.
	ListForwardTestsWorkflowParams struct{}

	// ListForwardTestsWorkflowResults is the output for the ListForwardTestsWorkflow.
	ListForwardTestsWorkflowResults struct {
		ForwardTests []forwardtest.ForwardTest
	}
)

// CreateForwardTestOrderWorkflowName is the name of the CreateForwardTestOrderWorkflow.
const CreateForwardTestOrderWorkflowName = "CreateForwardTestOrderWorkflow"

type (
	// CreateForwardTestOrderWorkflowParams is the input for the CreateForwardTestOrderWorkflow.
	CreateForwardTestOrderWorkflowParams struct {
		ForwardTestID uuid.UUID
		Order         order.Order
	}

	// CreateForwardTestOrderWorkflowResults is the output for the CreateForwardTestOrderWorkflow.
	CreateForwardTestOrderWorkflowResults struct{}
)

// ListForwardTestAccountsWorkflowName is the name of the ListForwardTestAccountsWorkflow.
const ListForwardTestAccountsWorkflowName = "ListForwardTestAccountsWorkflow"

type (
	// ListForwardTestAccountsWorkflowParams is the input for the ListForwardTestAccountsWorkflow.
	ListForwardTestAccountsWorkflowParams struct {
		ForwardTestID uuid.UUID
	}

	// ListForwardTestAccountsWorkflowResults is the output for the ListForwardTestAccountsWorkflow.
	ListForwardTestAccountsWorkflowResults struct {
		Accounts map[string]account.Account
	}
)

// GetForwardTestStatusWorkflowName is the name of the GetForwardTestStatusWorkflow.
const GetForwardTestStatusWorkflowName = "GetForwardTestStatusWorkflow"

type (
	// GetForwardTestStatusWorkflowParams is the input for the GetForwardTestStatusWorkflow.
	GetForwardTestStatusWorkflowParams struct {
		ForwardTestID uuid.UUID
	}

	// GetForwardTestStatusWorkflowResults is the output for the GetForwardTestStatusWorkflow.
	GetForwardTestStatusWorkflowResults struct {
		Status forwardtest.Status
	}
)
