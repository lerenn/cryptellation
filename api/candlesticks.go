package api

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
	// ListCandlesticksParams is the parameters of the ListCandlesticks activity.
	ListCandlesticksParams struct {
		Exchange string
		Pair     string
		Period   period.Symbol
		Start    *time.Time
		End      *time.Time
		Limit    uint
	}

	// ListCandlesticksResults is the result of the ListCandlesticks activity.
	ListCandlesticksResults struct {
		List *candlestick.List
	}
)
