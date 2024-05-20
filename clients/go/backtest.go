package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/event"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
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

		for {
			// Receiving events
			evt := <-events

			// Update the time of the run
			b.run.Time = evt.Time

			// If status, then there is no more events
			if evt.Type == event.TypeIsStatus {
				status := evt.Content.(event.Status)
				if status.Finished {
					endBacktest = true
				}

				break
			}

			// Call the bot on other events
			if err := b.bot.OnEvent(evt); err != nil {
				return err
			}
		}
	}

	return b.bot.OnExit()
}
