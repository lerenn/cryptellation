package temporal

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/models/timeserie"
)

const (
	// ListSMAWorkflowName is the name of the workflow to list SMA points.
	ListSMAWorkflowName = "ListSMAWorkflow"
)

type (
	// ListSMAWorkflowParams is the parameters of the ListSMA workflow.
	ListSMAWorkflowParams struct {
		Exchange     string
		Pair         string
		Period       period.Symbol
		Start        time.Time
		End          time.Time
		PeriodNumber int
		PriceType    candlestick.PriceType
	}

	// ListSMAWorkflowResults is the result of the ListSMA workflow.
	ListSMAWorkflowResults struct {
		Data *timeserie.TimeSerie[float64]
	}
)
