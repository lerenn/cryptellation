package api

import "github.com/lerenn/cryptellation/v1/pkg/models/exchange"

const (
	// ListExchangesWorkflowName is the name of the workflow to list exchanges.
	ListExchangesWorkflowName = "ListExchangesWorkflow"
)

type (
	// ListExchangesParams is the parameters of the ListExchanges activity.
	ListExchangesParams struct {
		Names []string
	}

	// ListExchangesResults is the result of the ListExchanges activity.
	ListExchangesResults struct {
		List []exchange.Exchange
	}
)
