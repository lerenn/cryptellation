package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/controllers/http"
)

type Daemon struct {
	// Application fields
	adapters    adapters
	components  components
	controllers controllers

	// Specific daemon fields
	health *http.Health
}

func New(ctx context.Context) (Daemon, error) {
	// Init and serve health server
	// NOTE: health OK, but not-ready yet
	h, err := http.NewHealth(ctx)
	if err != nil {
		return Daemon{}, err
	}

	// Init adapters
	adapters, err := newAdapters(ctx)
	if err != nil {
		return Daemon{}, err
	}

	// Init components
	components := newComponents(adapters)

	// Init controllers
	controllers, err := newControllers(components)
	if err != nil {
		adapters.Close(ctx)
		return Daemon{}, err
	}

	return Daemon{
		// Application specific
		adapters:    adapters,
		components:  components,
		controllers: controllers,
		// Daemon specific
		health: h,
	}, nil
}

func (d Daemon) Serve(ctx context.Context) error {
	if err := d.controllers.AsyncListen(); err != nil {
		return err
	}

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Serving health server
	go d.health.HTTPServe(ctx)

	// Service marked as ready
	telemetry.L(ctx).Info("Service is ready")
	d.health.Ready(true)

	// Wait for interrupt
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		telemetry.L(ctx).Info("Got SIGINT...")
	case syscall.SIGTERM:
		telemetry.L(ctx).Info("Got SIGTERM...")
	}

	telemetry.L(ctx).Info("The service is shutting down...")
	return nil
}

func (d Daemon) Close(ctx context.Context) {
	// Set daemon as not ready
	d.health.Ready(false)

	// Uninitialize application
	d.controllers.Close()
	d.components.Close()
	d.adapters.Close(ctx)
}
