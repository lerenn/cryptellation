package client

import (
	"context"

	"cryptellation/pkg/adapters/telemetry"
	"cryptellation/pkg/models/event"

	backtests "cryptellation/svc/backtests/clients/go"

	"cryptellation/svc/ticks/pkg/tick"
)

type Backtest struct {
	run *Run
	bot Bot
}

func NewBacktest(
	ctx context.Context,
	services Services,
	parameters backtests.BacktestCreationPayload,
	bot Bot,
) (*Backtest, error) {
	var run Run

	// Create the backtest
	id, err := services.Backtests().Create(ctx, parameters)
	if err != nil {
		return nil, err
	}

	// Set run
	run.ID = id
	run.Mode = ModeIsBacktest
	run.Services = services
	run.Time = parameters.StartTime

	// Init the robot
	bot.OnInit(ctx, &run)

	return &Backtest{
		run: &run,
		bot: bot,
	}, nil
}

func (b *Backtest) Run(ctx context.Context) error {
	// Listen to events
	events, err := b.run.Services.Backtests().ListenEvents(ctx, b.run.ID)
	if err != nil {
		return err
	}

	// Subscribe to ticks
	for _, ts := range b.bot.TicksToListen(ctx) {
		if err := b.run.Services.Backtests().Subscribe(ctx, b.run.ID, ts.Exchange, ts.Pair); err != nil {
			return err
		}
	}

	for endBacktest := false; !endBacktest; {
		err := b.run.Services.Backtests().Advance(ctx, b.run.ID)
		if err != nil {
			return err
		}
		telemetry.L(ctx).Debug("Backtest advanced")

		endBacktest, err = b.loopOnEvents(ctx, events)
		if err != nil {
			return err
		}
	}

	return b.bot.OnExit(ctx)
}

func (b *Backtest) loopOnEvents(ctx context.Context, events <-chan event.Event) (bool, error) {
	for {
		// Receiving events
		telemetry.L(ctx).Debug("Wait event")
		evt := <-events

		// Update the time of the run
		b.run.Time = evt.Time

		telemetry.L(ctx).Debugf("Event %q received", evt.Type.String())
		switch evt.Type {
		case event.TypeIsStatus:
			status := evt.Content.(event.Status)
			return status.Finished, nil // Exit loop event with indication wether the backtest is finished
		case event.TypeIsTick:
			t, ok := evt.Content.(tick.Tick)
			if !ok {
				telemetry.L(ctx).Error("tick event received but content is not a tick")
				continue
			}

			if err := b.bot.OnTick(ctx, t); err != nil {
				telemetry.L(ctx).Errorf("error on tick event: %w", err)
			}
		}
	}
}
