package main

import (
	"context"
	"fmt"

	cryptellation "github.com/lerenn/cryptellation/clients/go"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/event"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Bot struct {
	run          *cryptellation.Run
	isPositioned bool
}

func (b *Bot) OnInit(run *cryptellation.Run) {
	b.run = run
}

func (b *Bot) TicksToListen() []event.TickSubscription {
	return []event.TickSubscription{
		{
			Exchange: "binance",
			Pair:     "BTC-USDT",
		},
	}
}

func (b *Bot) OnEvent(event event.Event) error {
	t := event.Content.(tick.Tick)

	s, err := b.run.Services.Indicators().SMA(context.Background(), indicators.SMAPayload{
		Exchange:     t.Exchange,
		Pair:         t.Pair,
		Period:       period.M1,
		Start:        t.Time.Add(-period.M1.Duration()),
		End:          t.Time,
		PeriodNumber: 20,
		PriceType:    candlestick.PriceTypeIsClose,
	})
	if err != nil {
		return err
	}

	last, _ := s.Get(t.Time)
	previousLast, _ := s.Get(t.Time.Add(-period.M1.Duration()))
	if last > previousLast && !b.isPositioned {
		fmt.Println("+ at", t.Price)
		b.isPositioned = true

		if err := b.run.CreateOrder(context.Background(), backtests.OrderCreationPayload{
			Exchange: t.Exchange,
			Pair:     t.Pair,
			Quantity: 0.01,
			Side:     order.SideIsBuy,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}
	} else if last < previousLast && b.isPositioned {
		fmt.Println("- at", t.Price)
		b.isPositioned = false

		if err := b.run.CreateOrder(context.Background(), backtests.OrderCreationPayload{
			Exchange: t.Exchange,
			Pair:     t.Pair,
			Quantity: 0.01,
			Side:     order.SideIsSell,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) OnExit() error {
	if b.isPositioned {
		if err := b.run.CreateOrder(context.Background(), backtests.OrderCreationPayload{
			Exchange: "binance",
			Pair:     "BTC-USDT",
			Quantity: 0.01,
			Side:     order.SideIsSell,
			Type:     order.TypeIsMarket,
		}); err != nil {
			return err
		}
	}

	accounts, _ := b.run.GetAccounts(context.Background())
	fmt.Println("")
	fmt.Println("Result:", accounts["binance"].Balances["USDT"])

	return nil
}
