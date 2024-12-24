package binance

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	client "github.com/adshao/go-binance/v2"
	actutils "github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal/signals"
	"github.com/lerenn/cryptellation/v1/pkg/models/pair"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/temporal"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/activity"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Activities is the struct that will handle the exchanges.
type Activities struct {
	temporal temporalclient.Client
	*actutils.Binance
}

// New will create a new binance exchanges.
func New(temporal temporalclient.Client) (exchanges.Exchanges, error) {
	s, err := actutils.NewBinance(config.LoadBinanceTest())
	return &Activities{
		temporal: temporal,
		Binance:  s,
	}, err
}

// Name will return the name of the exchanges.
func (a *Activities) Name() string {
	return actutils.BinanceInfos.Name
}

// Register will register the exchanges.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.ListenSymbolActivity,
		activity.RegisterOptions{Name: exchanges.ListenSymbolActivityName})
}

// ListenSymbolActivity will listen to the symbol activity.
func (a *Activities) ListenSymbolActivity(
	ctx context.Context,
	params exchanges.ListenSymbolParams,
) (exchanges.ListenSymbolResults, error) {
	binanceSymbol, err := toBinanceSymbol(params.Symbol)
	if err != nil {
		return exchanges.ListenSymbolResults{}, err
	}

	// Start heartbeat on activity
	temporal.AsyncActivityHeartbeat(ctx, 300*time.Millisecond)

	// Listen to binance book ticker
	var lastBid, lastAsk string
	done, cancel, err := client.WsBookTickerServe(binanceSymbol, func(event *client.WsBookTickerEvent) {
		// Skip if same price as last tick
		if event.BestAskPrice == lastAsk && event.BestBidPrice == lastBid {
			return
		}
		lastAsk = event.BestAskPrice
		lastBid = event.BestBidPrice

		// Convert to tick
		t, err := toTick(params.Symbol, event.BestAskPrice, event.BestBidPrice)
		if err != nil {
			telemetry.L(ctx).Errorf("Failed to convert binance tick: %v", err)
			return
		}

		// Send it to main workflow through Signal
		err = a.temporal.SignalWorkflow(ctx, params.ParentWorkflowID, "", signals.NewTickReceivedSignalName, t)
		if err != nil {
			a.handleNewTickSignalError(ctx, err, params)
		}
	}, nil)
	if err != nil {
		return exchanges.ListenSymbolResults{}, err
	}

	// Wait for context to be done or cancelled
	select {
	case <-done:
		// If done, return error as listener stopped
		return exchanges.ListenSymbolResults{}, fmt.Errorf("binance listener stopped")
	case <-ctx.Done():
		// If context is done, cancel listener and return
		cancel <- struct{}{}
		return exchanges.ListenSymbolResults{}, nil
	}
}

func (a *Activities) handleNewTickSignalError(ctx context.Context, ntsErr error, params exchanges.ListenSymbolParams) {
	// Context was cancelled, this will stop listener
	if errors.Is(ntsErr, context.Canceled) {
		return
	}

	// Check if parent workflow is still running
	desc, err := a.temporal.DescribeWorkflowExecution(ctx, params.ParentWorkflowID, "")
	if err != nil {
		telemetry.L(ctx).Errorf("Failed to describe parent workflow: %v", err)
		return
	} else if desc.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_COMPLETED {
		// Workflow is already completed
		return
	}

	// That shouldn't happen, log error
	telemetry.L(ctx).Errorf("Failed to signal binance tick: %v", err)
}

func toTick(symbol, ask, bid string) (tick.Tick, error) {
	askPrice, err := strconv.ParseFloat(ask, 64)
	if err != nil {
		return tick.Tick{}, err
	}

	bidPrice, err := strconv.ParseFloat(bid, 64)
	if err != nil {
		return tick.Tick{}, err
	}

	return tick.Tick{
		Time:     time.Now().UTC(),
		Exchange: "binance",
		Pair:     symbol,
		Price:    (askPrice + bidPrice) / 2,
	}, nil
}

func toBinanceSymbol(symbol string) (string, error) {
	base, quote, err := pair.ParsePair(symbol)
	return fmt.Sprintf("%s%s", base, quote), err
}
