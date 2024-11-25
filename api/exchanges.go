package api

import "github.com/lerenn/cryptellation/v1/pkg/models/exchange"

const (
	// GetExchangeWorkflowName is the name of the workflow to get an exchange.
	GetExchangeWorkflowName = "GetExchangeWorkflow"
)

type (
	// GetExchangeParams is the parameters of the GetExchange activity.
	GetExchangeParams struct {
		Name string
	}

	// GetExchangeResults is the result of the GetExchange activity.
	GetExchangeResults struct {
		Exchange exchange.Exchange
	}
)

const (
	// ListExchangesWorkflowName is the name of the workflow to list exchanges.
	ListExchangesWorkflowName = "ListExchangesWorkflow"
)

type (
	// ListExchangesParams is the parameters of the ListExchanges activity.
	ListExchangesParams struct{}

	// ListExchangesResults is the result of the ListExchanges activity.
	ListExchangesResults struct {
		List []string
	}
)
