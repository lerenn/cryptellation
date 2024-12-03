package api

import "github.com/lerenn/cryptellation/v1/pkg/models/exchange"

const (
	// GetExchangeWorkflowName is the name of the workflow to get an exchange.
	GetExchangeWorkflowName = "GetExchangeWorkflow"
)

type (
	// GetExchangeWorkflowParams is the parameters of the GetExchange workflow.
	GetExchangeWorkflowParams struct {
		Name string
	}

	// GetExchangeWorkflowResults is the result of the GetExchange workflow.
	GetExchangeWorkflowResults struct {
		Exchange exchange.Exchange
	}
)

const (
	// ListExchangesWorkflowName is the name of the workflow to list exchanges.
	ListExchangesWorkflowName = "ListExchangesWorkflow"
)

type (
	// ListExchangesWorkflowParams is the parameters of the ListExchanges workflow.
	ListExchangesWorkflowParams struct{}

	// ListExchangesWorkflowResults is the result of the ListExchanges workflow.
	ListExchangesWorkflowResults struct {
		List []string
	}
)
