package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/event"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

type Backtest struct {
	run *Run
	bot Bot
}

func NewBacktest(services Services, parameters backtests.BacktestCreationPayload, bot Bot) (*Backtest, error) {
	var run Run

	// Create the backtest
	id, err := services.Backtests().Create(context.Background(), parameters)
	if err != nil {
		return nil, err
	}

	// Set run
	run.ID = id
	run.Mode = ModeIsBacktest
	run.Services = services
	run.Time = parameters.StartTime

	// Init the robot
	bot.OnInit(&run)

	return &Backtest{
		run: &run,
		bot: bot,
	}, nil
}

func (b *Backtest) Run() error {
	// Listen to events
	events, err := b.run.Services.Backtests().ListenEvents(context.Background(), b.run.ID)
	if err != nil {
		return err
	}

	// Subscribe to ticks
	for _, ts := range b.bot.TicksToListen() {
		if err := b.run.Services.Backtests().Subscribe(context.Background(), b.run.ID, ts.Exchange, ts.Pair); err != nil {
			return err
		}
	}

	for endBacktest := false; !endBacktest; {
		err := b.run.Services.Backtests().Advance(context.Background(), b.run.ID)
		if err != nil {
			return err
		}
		telemetry.L(context.Background()).Debug("Backtest advanced")

		endBacktest, err = b.loopOnEvents(events)
		if err != nil {
			return err
		}
	}

	return b.bot.OnExit()
}

func (b *Backtest) loopOnEvents(events <-chan event.Event) (bool, error) {
	for {
		// Receiving events
		telemetry.L(context.Background()).Debug("Wait event")
		evt := <-events

		// Update the time of the run
		b.run.Time = evt.Time

		telemetry.L(context.Background()).Debugf("Event %q received", evt.Type.String())
		switch evt.Type {
		case event.TypeIsStatus:
			status := evt.Content.(event.Status)
			return status.Finished, nil // Exit loop event with indication wether the backtest is finished
		case event.TypeIsTick:
			t, ok := evt.Content.(tick.Tick)
			if !ok {
				telemetry.L(context.Background()).Error("tick event received but content is not a tick")
				continue
			}

			if err := b.bot.OnTick(t); err != nil {
				return false, err
			}
		}
	}
}
