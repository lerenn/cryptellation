package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/activities/exchanges"
	"github.com/lerenn/cryptellation/v1/pkg/domains/ticks/internal"
	"github.com/lerenn/cryptellation/v1/pkg/models/pair"
	"github.com/lerenn/cryptellation/v1/pkg/models/tick"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	temporalclient "go.temporal.io/sdk/client"
)

type Activities struct {
	temporal temporalclient.Client
	*activities.Binance
}

func New(temporal temporalclient.Client) (*Activities, error) {
	s, err := activities.NewBinance(config.LoadBinanceTest())
	return &Activities{
		temporal: temporal,
		Binance:  s,
	}, err
}

func (s *Activities) ListenSymbol(ctx context.Context, params exchanges.ListenSymbolParams) (exchanges.ListenSymbolResults, error) {
	binanceSymbol, err := toBinanceSymbol(params.Symbol)
	if err != nil {
		return exchanges.ListenSymbolResults{}, err
	}

	var lastBid, lastAsk string
	_, _, err = client.WsBookTickerServe(binanceSymbol, func(event *client.WsBookTickerEvent) {
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
		err = s.temporal.SignalWorkflow(ctx, params.ParentWorkflowID, "", internal.NewTickReceivedSignalName, t)
		if err != nil {
			telemetry.L(ctx).Errorf("Failed to signal binance tick: %v", err)
			return
		}
	}, nil)

	// TODO: manage when error or done

	return exchanges.ListenSymbolResults{}, err
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
