package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/internal/adapters/db/sql"
	"github.com/digital-feather/cryptellation/internal/adapters/exchanges"
	"github.com/digital-feather/cryptellation/internal/components/candlesticks"
	"github.com/digital-feather/cryptellation/internal/controllers/nats"
	natsCandlesticks "github.com/digital-feather/cryptellation/internal/controllers/nats/candlesticks"
	"github.com/digital-feather/cryptellation/pkg/http/health"
)

func initComponent() (candlesticks.Port, error) {
	// Init database client
	db, err := sql.New(sql.LoadConfigFromEnv())
	if err != nil {
		return nil, err
	}

	// Init exchanges connections
	exchanges, err := exchanges.New(exchanges.LoadConfigFromEnv())
	if err != nil {
		return nil, err
	}

	// Init component
	return candlesticks.New(db, exchanges), nil
}

func initController(component candlesticks.Port) (func(), error) {
	// Init NATS controller
	natsController, err := natsCandlesticks.NewServer(nats.LoadConfigFromEnv(), component)
	if err != nil {
		return func() {}, err
	}

	// Listen on NATS controller
	if err := natsController.Listen(); err != nil {
		return func() {}, err
	}

	return func() {
		natsController.Close()
	}, nil
}

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init component
	component, err := initComponent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initialize application: %w", err))
		return 255
	}

	// Init controller
	cleanupController, err := initController(component)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initializing NATS controller: %w", err))
		return 255
	}
	defer cleanupController()

	// Service marked as ready
	log.Println("Service is ready")
	h.Ready(true)

	// Wait for interrupt
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	return 0
}

func main() {
	os.Exit(run())
}
