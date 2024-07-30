package client

import (
	"context"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"cryptellation/svc/ticks/pkg/tick"
)

type ForwardTest struct {
	run *Run
	bot Bot
}

func NewForwardTest(services Services, parameters forwardtest.NewPayload, bot Bot) (*ForwardTest, error) {
	var run Run

	// Create the forward test
	id, err := services.ForwardTests().CreateForwardTest(context.Background(), parameters)
	if err != nil {
		return nil, err
	}

	// Set run
	run.ID = id
	run.Mode = ModeIsForwardTest
	run.Services = services

	// Init the robot
	bot.OnInit(&run)

	return &ForwardTest{
		run: &run,
		bot: bot,
	}, nil
}

func (ft *ForwardTest) Run() error {
	// Subscribe to ticks
	tListen := ft.bot.TicksToListen()
	ticksChan := make(chan tick.Tick, 64)
	for _, ts := range tListen {
		ctx, cancel := context.WithCancel(context.Background())
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
		if err := ft.bot.OnTick(t); err != nil {
			return err
		}
	}

	return ft.bot.OnExit()
}
