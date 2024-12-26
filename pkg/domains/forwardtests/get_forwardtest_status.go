package forwardtests

import (
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
	"github.com/lerenn/cryptellation/v1/pkg/utils"
	"go.temporal.io/sdk/workflow"
)

var (
	// ErrNoActualPrice is the error when there is no actual price when requesting status.
	ErrNoActualPrice = fmt.Errorf("no actual price")
)

const (
	// DefaultBalanceSymbol is the default symbol used to have the total balance.
	DefaultBalanceSymbol = "USDT"
)

// GetForwardTestStatusWorkflow is the workflow to get the forwardtest status.
func (wf *workflows) GetForwardTestStatusWorkflow(
	ctx workflow.Context,
	params api.GetForwardTestStatusWorkflowParams,
) (api.GetForwardTestStatusWorkflowResults, error) {
	// Read forwardtest from database
	ft, err := wf.readForwardTestFromDB(ctx, params.ForwardTestID)
	if err != nil {
		return api.GetForwardTestStatusWorkflowResults{},
			fmt.Errorf("could not read forwardtest from db: %w", err)
	}

	// Get value for each symbol in accounts
	total := 0.0
	for exchange, account := range ft.Accounts {
		for symbol, balance := range account.Balances {
			if symbol == DefaultBalanceSymbol {
				total += balance
				continue
			}

			// Get price
			p := symbol + "-" + DefaultBalanceSymbol
			csRes, err := wf.cryptellation.ListCandlesticks(ctx, api.ListCandlesticksWorkflowParams{
				Exchange: exchange,
				Pair:     p,
				Period:   period.M1,
				Start:    utils.ToReference(time.Now().Add(-time.Minute * 10)),
				End:      utils.ToReference(time.Now()),
				Limit:    1,
			}, nil)
			if err != nil {
				return api.GetForwardTestStatusWorkflowResults{},
					fmt.Errorf("could not get candlesticks from service: %w", err)
			}

			c, ok := csRes.List.Last()
			if !ok {
				return api.GetForwardTestStatusWorkflowResults{}, fmt.Errorf("%w: %s", ErrNoActualPrice, p)
			}

			// Calculate value
			total += balance * c.Close
		}
	}

	return api.GetForwardTestStatusWorkflowResults{
		Status: forwardtest.Status{
			Balance: total,
		},
	}, nil
}
