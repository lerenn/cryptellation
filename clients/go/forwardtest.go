package client

import (
	"context"

	"github.com/lerenn/cryptellation/forwardtests/pkg/forwardtest"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"

	"github.com/lerenn/cryptellation/ticks/pkg/tick"
)

type ForwardTest struct {
	run *Run
	bot Bot
}

func NewForwardTest(
	ctx context.Context,
	services Services,
	parameters forwardtest.NewPayload,
	bot Bot,
) (*ForwardTest, error) {
	var run Run

	// Create the forward test
	id, err := services.ForwardTests().CreateForwardTest(ctx, parameters)
	if err != nil {
		return nil, err
	}

	// Set run
	run.ID = id
	run.Mode = ModeIsForwardTest
	run.Services = services

	// Init the robot
	bot.OnInit(ctx, &run)

	return &ForwardTest{
		run: &run,
		bot: bot,
	}, nil
}

func (ft *ForwardTest) Run(ctx context.Context) error {
	// Subscribe to ticks
	tListen := ft.bot.TicksToListen(ctx)
	ticksChan := make(chan tick.Tick, 64)
	for _, ts := range tListen {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch, err := ft.run.Services.Ticks().SubscribeToTicks(ctx, ts)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case t := <-ch:
					ticksChan <- t
				}
			}
		}()
	}

	for t := range ticksChan {
		// Update the time of the run
		ft.run.Time = t.Time

		// Call the bot on other events
		if err := ft.bot.OnTick(ctx, t); err != nil {
			telemetry.L(ctx).Error("error on tick: " + err.Error())
			continue
		}
	}

	return ft.bot.OnExit(ctx)
}
