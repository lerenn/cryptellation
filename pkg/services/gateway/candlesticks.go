package gateway

import (
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/lerenn/cryptellation/v1/api/worker/go"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

func getCandlesticksParamsFromAPIToWorker(params GetCandlesticksParams) (api.ListCandlesticksWorkflowParams, error) {
	// Get and check the period
	per, err := period.FromString(params.Interval)
	if err != nil {
		return api.ListCandlesticksWorkflowParams{}, err
	}

	// Get and check the start time
	var start *time.Time
	if params.StartTime != nil {
		s, err := time.Parse(time.RFC3339, *params.StartTime)
		if err != nil {
			return api.ListCandlesticksWorkflowParams{}, err
		}
		start = &s
	}

	// Get and check the end time
	var end *time.Time
	if params.EndTime != nil {
		e, err := time.Parse(time.RFC3339, *params.EndTime)
		if err != nil {
			return api.ListCandlesticksWorkflowParams{}, err
		}
		end = &e
	}

	return api.ListCandlesticksWorkflowParams{
		Exchange: params.Exchange,
		Pair:     params.Symbol,
		Period:   per,
		Start:    start,
		End:      end,
	}, nil
}

// GetCandlesticks is the handler for the /candlesticks endpoint.
// It returns the candlesticks for the given parameters.
func (s *Server) GetCandlesticks(c *gin.Context, params GetCandlesticksParams) {
	// Get the parameters from the request
	apiParams, err := getCandlesticksParamsFromAPIToWorker(params)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// List candlesticks
	res, err := s.client.ListCandlesticks(c.Request.Context(), apiParams)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert candlesticks to the API format
	candlesticks := make([]Candlestick, 0, res.List.Data.Len())
	_ = res.List.Data.Loop(func(t time.Time, d candlestick.Candlestick) (bool, error) {
		candlesticks = append(candlesticks, Candlestick{
			Time:       t.Format(time.RFC3339),
			Close:      float32(d.Close),
			High:       float32(d.High),
			Low:        float32(d.Low),
			Open:       float32(d.Open),
			Volume:     float32(d.Volume),
			Uncomplete: nil,
		})
		return false, nil
	})

	// Return the candlesticks
	c.JSON(200, candlesticks)
}
