package temporal

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

const (
	// ListCandlesticksWorkflowName is the name of the workflow to get candlesticks.
	ListCandlesticksWorkflowName = "ListCandlesticksWorkflow"
)

type (
	// ListCandlesticksWorkflowParams is the parameters of the ListCandlesticks workflow.
	ListCandlesticksWorkflowParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    *time.Time
		End      *time.Time
		Limit    uint
	}

	// ListCandlesticksWorkflowResults is the result of the ListCandlesticks workflow.
	ListCandlesticksWorkflowResults struct {
		List *candlestick.List
	}
)
