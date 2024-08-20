package bot

import (
	"context"
	"fmt"
	"time"

	cryptellation "github.com/lerenn/cryptellation/clients/go"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Bot struct {
	run          *cryptellation.Run
	isPositioned bool
}

func (b *Bot) OnInit(ctx context.Context, run *cryptellation.Run) {
	telemetry.L(ctx).Debug("OnInit() called")
	defer telemetry.L(ctx).Info("Bot initialized")

	b.run = run
}

func (b *Bot) TicksToListen(ctx context.Context) []event.TickSubscription {
	telemetry.L(ctx).Debug("TicksToListen() called")

	return []event.TickSubscription{
		{
			Exchange: "binance",
			Pair:     "BTC-USDT",
		},
	}
}

func (b *Bot) OnTick(ctx context.Context, t tick.Tick) error {
	telemetry.L(ctx).Debugf("OnTick(t=%+v) called", t)

	payload := indicators.SMAPayload{
		Exchange:     t.Exchange,
		Pair:         t.Pair,
		Period:       period.M1,
		Start:        t.Time.Add(-period.M1.Duration() * 2),
		End:          t.Time.Add(-period.M1.Duration()),
		PeriodNumber: 20,
		PriceType:    candlestick.PriceTypeIsClose,
	}
	telemetry.L(ctx).Debugf("Request SMA with %+v", payload)
	s, err := b.run.Services.Indicators().SMA(ctx, payload)
	if err != nil {
		return err
	}

	_ = s.Loop(func(t time.Time, v float64) (bool, error) {
		telemetry.L(ctx).Debugf("SMA point at %s: %f", t, v)
		return false, nil
	})

	tLast, last, ok := s.Last()
	if !ok {
		return fmt.Errorf("last SMA is empty")
	}
	previousLast, ok := s.Get(tLast.Add(-period.M1.Duration()))
	if !ok {
		return fmt.Errorf("previous SMA is empty")
	}
	telemetry.L(ctx).Debugf("SMA received [%f, %f]", previousLast, last)

	if last > previousLast && !b.isPositioned {
		fmt.Println("+ at", t.Price)
		b.isPositioned = true

		if err := b.run.CreateOrder(ctx, common.OrderCreationPayload{
			Exchange: t.Exchange,
			Pair:     t.Pair,
			Quantity: 0.01,
			Side:     order.SideIsBuy,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}

		telemetry.L(ctx).Debug("Buy order created")
	} else if last < previousLast && b.isPositioned {
		fmt.Println("- at", t.Price)
		b.isPositioned = false

		if err := b.run.CreateOrder(ctx, common.OrderCreationPayload{
			Exchange: t.Exchange,
			Pair:     t.Pair,
			Quantity: 0.01,
			Side:     order.SideIsSell,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}

		telemetry.L(ctx).Debug("Sell order created")
	} else {
		telemetry.L(ctx).Debug("No action taken")
	}

	return nil
}

func (b *Bot) OnExit(ctx context.Context) error {
	telemetry.L(ctx).Debug("OnExit() called")

	if b.isPositioned {
		if err := b.run.CreateOrder(ctx, common.OrderCreationPayload{
			Exchange: "binance",
			Pair:     "BTC-USDT",
			Quantity: 0.01,
			Side:     order.SideIsSell,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}
	}

	accounts, err := b.run.GetAccounts(ctx)
	if err != nil {
		return err
	}
	telemetry.L(ctx).Infof("Result: $%f", accounts["binance"].Balances["USDT"])

	return nil
}
